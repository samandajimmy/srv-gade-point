package http_test

import (
	"encoding/json"
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	rwdhttp "gade/srv-gade-point/rewards/delivery/http"
	"gade/srv-gade-point/test"
	"gade/srv-gade-point/test/fakedata"
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	db       *database.DbConfig
	migrator *migrate.Migrate
	usecases config.Usecases
	e        test.DummyEcho

	_ = BeforeSuite(func() {
		config.LoadEnv()
		config.LoadTestData()
		db, migrator = test.NewDbConn()
	})

	_ = AfterSuite(func() {
		_ = migrator.Drop()

		db.Sql.Close()
		migrator.Close()
	})
)

var _ = Describe("Reward", func() {
	var pl interface{}
	var response models.Response
	var responseData map[string]interface{}
	var expectResp models.Response
	var handler rwdhttp.RewardHandler
	var campaign models.Campaign
	var reward models.Reward
	var rewards []models.Reward
	var voucher models.Voucher
	var refCode models.RespReferral
	var createReferral models.RequestCreateReferral

	BeforeEach(func() {
		response, expectResp, reward, voucher = models.Response{}, models.Response{}, models.Reward{}, models.Voucher{}
		rewards = []models.Reward{}
		e = test.NewDummyEcho(http.MethodPost, "/", pl)
		_, usecases = test.LoadRealRepoUsecase(db)
		handler.RewardUseCase = usecases.RewardUseCase
		campaign = fakedata.Campaign()
	})

	Describe("MultiRewardInquiry", func() {
		var reqpl models.PayloadValidator
		var expectedValue float64

		JustBeforeEach(func() {
			rewards = append(rewards, reward)
			campaign.Rewards = &rewards
			refCode = createCampaign(&campaign, createReferral)
			reqpl = mapReqpl(reqpl, reward, voucher, refCode)

			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.MultiRewardInquiry(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		BeforeEach(func() {
			reqpl = fakedata.PayloadInquiry()
		})

		Context("inquiry succeeded", func() {
			JustBeforeEach(func() {
				responseData = response.Data.(map[string]interface{})
				expectResp = models.Response{
					Code:    "00",
					Status:  "Success",
					Message: "Data Berhasil Dikirim",
				}
				expectResp.Data = inqResponse(responseData, reward, expectedValue, reqpl.CIF)
			})

			Context("types of reward", func() {
				Context("inquiry reward direct discount", func() {
					BeforeEach(func() {
						reward = fakedata.RewardDirectDisc(true)
						expectedValue = models.RoundDown(*reward.Validators.Value, 0)
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("inquiry reward percent discount", func() {
					Context("with max value", func() {
						BeforeEach(func() {
							reward = fakedata.RewardPercentDisc(true, true)
							reward.Validators.Discount = helper.CreateFloat64(10)
							reward.Validators.MaxValue = helper.CreateFloat64(5000)
						})

						Context("reward greater than and equals to max value", func() {
							BeforeEach(func() {
								reqpl.TransactionAmount = helper.CreateFloat64(60000)
								expectedValue = *reward.Validators.MaxValue
							})

							Context("reward greater than", func() {
								It("expect response to return as expected", func() {
									Expect(response).To(Equal(expectResp))
								})
							})

							Context("reward equals", func() {
								It("expect response to return as expected", func() {
									Expect(response).To(Equal(expectResp))
								})
							})
						})

						Context("reward less than max value", func() {
							BeforeEach(func() {
								reqpl.TransactionAmount = helper.CreateFloat64(6000)
								expectedValue = 600
							})

							It("expect response to return as expected", func() {
								Expect(response).To(Equal(expectResp))
							})
						})
					})

					Context("without max value", func() {
						BeforeEach(func() {
							reward = fakedata.RewardPercentDisc(true, false)
							reward.Validators.Discount = helper.CreateFloat64(10)
							reqpl.TransactionAmount = helper.CreateFloat64(6000)
							expectedValue = *reqpl.TransactionAmount * (*reward.Validators.Discount / 100)
						})

						It("expect response to return as expected", func() {
							Expect(response).To(Equal(expectResp))
						})
					})

				})

				Context("inquiry reward goldback", func() {
					BeforeEach(func() {
						reward = fakedata.RewardGoldback(true)
						expectedValue = *reqpl.TransactionAmount / (*reward.Validators.Deviden)
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("inquiry reward voucher", func() {
					BeforeEach(func() {
						voucher = fakedata.VoucherDirectDisc()
						_ = usecases.VoucherUseCase.CreateVoucher(e.Context, &voucher)

						reward = fakedata.RewardVoucher(true, voucher.ID)
					})

					JustBeforeEach(func() {
						expectResp.Data = map[string]interface{}{
							"refTrx": responseData["refTrx"],
							"rewards": []interface{}{
								map[string]interface{}{
									"product":        reward.Validators.Product,
									"refTrx":         responseData["refTrx"],
									"type":           models.RewardTypeStr[*reward.Type],
									"journalAccount": reward.JournalAccount,
									"voucherName":    voucher.Name,
									"cif":            reqpl.CIF,
								},
							},
						}
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("inquiry reward incentive", func() {
					BeforeEach(func() {
						reward = fakedata.RewardIncentive(true)
						expectedValue = 0
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

			})

			Context("types of promo code", func() {
				Context("promo", func() {
					BeforeEach(func() {
						reward = fakedata.RewardDirectDisc(true)
						expectedValue = models.RoundDown(*reward.Validators.Value, 0)
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("voucher", func() {
					BeforeEach(func() {
						voucher = fakedata.VoucherDirectDisc()
						_ = usecases.VoucherUseCase.CreateVoucher(e.Context, &voucher)
						payload := map[string]interface{}{
							"voucherId": strconv.FormatInt(voucher.ID, 10),
							"page":      0,
							"limit":     0,
						}
						voucherCodes, _, _ := usecases.VoucherCodeUseCase.GetVoucherCodes(e.Context, payload)
						voucher.PromoCode = &voucherCodes[0].PromoCode

						expectedValue = models.RoundDown(*voucher.Validators.Value, 0)
					})

					JustBeforeEach(func() {
						expectResp.Data = map[string]interface{}{
							"refTrx": responseData["refTrx"],
							"rewards": []interface{}{
								map[string]interface{}{
									"product":        voucher.Validators.Product,
									"refTrx":         responseData["refTrx"],
									"type":           models.RewardTypeStr[*voucher.Type],
									"journalAccount": voucher.JournalAccount,
									"value":          expectedValue,
									"cif":            reqpl.CIF,
								},
							},
						}
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("referral", func() {
					var reward1 models.Reward

					BeforeEach(func() {
						campaign = fakedata.CampaignReferral()
						reward = fakedata.RewardDirectDisc(false)
						reward1 = fakedata.RewardIncentive(false)
						reward1.Validators.Channel = reward.Validators.Channel
						reward1.Validators.TransactionType = reward.Validators.TransactionType
						reward1.Validators.Product = reward.Validators.Product
						rewards = append(rewards, reward1)

						createReferral = models.RequestCreateReferral{
							CIF:    gofakeit.Regex("[1234567890]{10}"),
							Prefix: campaign.Metadata.Prefix,
						}
					})

					JustBeforeEach(func() {
						rwds := responseData["rewards"].([]interface{})
						refTrx1 := rwds[0].(map[string]interface{})["refTrx"].(string)
						refTrx2 := rwds[1].(map[string]interface{})["refTrx"].(string)
						rwd1 := rwdInqResponse(refTrx2, reward.Validators.Product, reqpl.CIF, reward.JournalAccount, *reward.Type, models.RoundDown(*reward.Validators.Value, 0))
						rwd2 := rwdInqResponse(refTrx1, reward1.Validators.Product, refCode.CIF, reward1.JournalAccount, *reward1.Type, 0)
						rwd2["reference"] = reward1.Validators.Target
						rwd1["reference"] = "referral"
						expectResp.Data = map[string]interface{}{
							"refTrx":  responseData["refTrx"],
							"rewards": []interface{}{rwd2, rwd1},
						}
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})
			})

			Context("multiple rewards", func() {
				var reward1, reward2 models.Reward
				var rwd, rwd1, rwd2 map[string]interface{}

				BeforeEach(func() {
					reward = fakedata.RewardDirectDisc(true)
					reward1 = fakedata.RewardPercentDisc(true, false)
					reward1.Validators.Channel = reward.Validators.Channel
					reward1.Validators.TransactionType = reward.Validators.TransactionType
					reward1.Validators.Product = reward.Validators.Product
					reward1.PromoCode = reward.PromoCode
					rewards = append(rewards, reward1)
					reward2 = fakedata.RewardGoldback(true)
					reward2.Validators.Channel = reward.Validators.Channel
					reward2.Validators.TransactionType = reward.Validators.TransactionType
					reward2.Validators.Product = reward.Validators.Product
					reward2.PromoCode = reward.PromoCode
					rewards = append(rewards, reward2)
				})

				JustBeforeEach(func() {
					rwds := responseData["rewards"].([]interface{})
					refTrx1 := rwds[2].(map[string]interface{})["refTrx"].(string)
					refTrx2 := rwds[0].(map[string]interface{})["refTrx"].(string)
					refTrx3 := rwds[1].(map[string]interface{})["refTrx"].(string)
					rwd = rwdInqResponse(refTrx1, reward.Validators.Product, reqpl.CIF, reward.JournalAccount, *reward.Type, models.RoundDown(*reward.Validators.Value, 0))
					rwd1 = rwdInqResponse(refTrx2, reward1.Validators.Product, reqpl.CIF, reward1.JournalAccount, *reward1.Type, models.RoundDown(*reqpl.TransactionAmount*(*reward1.Validators.Discount/100), 0))
					rwd2 = rwdInqResponse(refTrx3, reward2.Validators.Product, reqpl.CIF, reward2.JournalAccount, *reward2.Type, *reqpl.TransactionAmount/(*reward2.Validators.Deviden))
					expectResp.Data = map[string]interface{}{
						"refTrx":  responseData["refTrx"],
						"rewards": []interface{}{rwd, rwd1, rwd2},
					}
				})

				It("expect response to return as expected", func() {
					Expect(responseData["rewards"]).To(ConsistOf([]interface{}{rwd, rwd1, rwd2}))
				})
			})

		})
	})

})

func inqResponse(responseData map[string]interface{}, reward models.Reward, value float64, cif string) map[string]interface{} {
	if reward == (models.Reward{}) {
		return map[string]interface{}{}
	}

	rwdInq := rwdInqResponse(responseData["refTrx"].(string), reward.Validators.Product, cif,
		reward.JournalAccount, *reward.Type, value)

	return map[string]interface{}{
		"refTrx":  responseData["refTrx"],
		"rewards": []interface{}{rwdInq},
	}
}

func rwdInqResponse(refTrx, product, cif, journalAccount string, rewardType int64, value float64) map[string]interface{} {
	rwdInq := map[string]interface{}{
		"product":        product,
		"refTrx":         refTrx,
		"type":           models.RewardTypeStr[rewardType],
		"journalAccount": journalAccount,
		"cif":            cif,
	}

	if value > 0 {
		rwdInq["value"] = value
	}

	return rwdInq
}

func mapReqpl(reqpl models.PayloadValidator, reward models.Reward, voucher models.Voucher, refCode models.RespReferral) models.PayloadValidator {
	if refCode != (models.RespReferral{}) {
		reqpl.Validators.Channel = reward.Validators.Channel
		reqpl.Validators.TransactionType = reward.Validators.TransactionType
		reqpl.Validators.Product = reward.Validators.Product
		reqpl.PromoCode = refCode.ReferralCode

		return reqpl
	}

	if reward != (models.Reward{}) {
		reqpl.Validators.Channel = reward.Validators.Channel
		reqpl.Validators.TransactionType = reward.Validators.TransactionType
		reqpl.Validators.Product = reward.Validators.Product
		reqpl.PromoCode = reward.PromoCode

		return reqpl
	}

	if voucher != (models.Voucher{}) {
		reqpl.Validators.Channel = voucher.Validators.Channel
		reqpl.Validators.TransactionType = voucher.Validators.TransactionType
		reqpl.Validators.Product = voucher.Validators.Product
		reqpl.PromoCode = *voucher.PromoCode

		return reqpl
	}

	return reqpl
}

func createCampaign(campaign *models.Campaign, createReferral models.RequestCreateReferral) models.RespReferral {
	_ = usecases.CampaignUseCase.CreateCampaign(e.Context, campaign)

	if createReferral == (models.RequestCreateReferral{}) {
		return models.RespReferral{}
	}

	rc, _ := usecases.ReferralsUseCase.UCreateReferralCodes(e.Context, createReferral)

	return rc
}
