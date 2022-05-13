package http_test

import (
	"encoding/json"
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	rwdhttp "gade/srv-gade-point/rewards/delivery/http"
	"gade/srv-gade-point/test"
	"gade/srv-gade-point/test/fakedata"
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	db             *database.DbConfig
	migrator       *migrate.Migrate
	pl             interface{}
	response       models.Response
	responseData   map[string]interface{}
	expectResp     models.Response
	handler        rwdhttp.RewardHandler
	campaign       models.Campaign
	reward         models.Reward
	rewards        []models.Reward
	voucher        models.Voucher
	refCode        models.RespReferral
	createReferral models.RequestCreateReferral
	usecases       config.Usecases
	repos          config.Repositories
	e              test.DummyEcho
	plPayment      models.RewardPayment
	withEcho       bool
	expectedValue  float64

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

	_ = BeforeEach(func() {
		pl = nil
		response = models.Response{}
		responseData = map[string]interface{}{}
		expectResp = models.Response{}
		handler = rwdhttp.RewardHandler{}
		campaign = models.Campaign{}
		reward = models.Reward{}
		rewards = []models.Reward{}
		voucher = models.Voucher{}
		refCode = models.RespReferral{}
		createReferral = models.RequestCreateReferral{}
		usecases = config.Usecases{}
		repos = config.Repositories{}
		e = test.DummyEcho{}
		withEcho = true
		expectedValue = 0
	})
)

var _ = Describe("Reward", func() {

	BeforeEach(func() {
		response, expectResp, response.Data = models.Response{}, models.Response{}, map[string]interface{}{}
		e, handler, reward, voucher = test.DummyEcho{}, rwdhttp.RewardHandler{}, models.Reward{}, models.Voucher{}
		rewards = []models.Reward{}
		withEcho = true

		campaign = fakedata.Campaign()
		repos = config.NewRepositories(db.Sql, db.Bun)
		usecases = config.NewUsecases(repos)
		handler.RewardUseCase = usecases.RewardUseCase
	})

	Describe("MultiRewardInquiry", func() {
		var reqpl models.PayloadValidator

		JustBeforeEach(func() {
			if !withEcho {
				return
			}

			rewards = append(rewards, reward)
			campaign.Rewards = &rewards
			refCode = createCampaign(&campaign, createReferral)
			reqpl = mapReqpl(reqpl, reward, voucher, refCode)
			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.MultiRewardInquiry(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
			newExpectedResp(reqpl.CIF)
		})

		BeforeEach(func() {
			reqpl = fakedata.PayloadInquiry()
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
					tempEcho := test.NewDummyEcho(http.MethodPost, "/", nil)
					_ = usecases.VoucherUseCase.CreateVoucher(tempEcho.Context, &voucher)

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

				Context("generic promo", func() {
					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("private promo", func() {
					Context("isPrivate = 1", func() {
						BeforeEach(func() {
							reward.Validators.IsPrivate = 1
						})

						It("expect response to return as expected", func() {
							Expect(response).To(Equal(expectResp))
						})
					})

					Context("isPrivate = 0", func() {
						BeforeEach(func() {
							reward.Validators.IsPrivate = 2
						})

						It("expect response to return as expected", func() {
							Expect(response).To(Equal(expectResp))
						})
					})
				})
			})

			Context("voucher", func() {
				var tempEcho test.DummyEcho

				BeforeEach(func() {
					tempEcho = test.NewDummyEcho(http.MethodPost, "/", nil)
					voucher = fakedata.VoucherDirectDisc()
					_ = usecases.VoucherUseCase.CreateVoucher(tempEcho.Context, &voucher)
					payload := map[string]interface{}{
						"voucherId": strconv.FormatInt(voucher.ID, 10),
						"page":      0,
						"limit":     0,
					}

					voucherCodes, _, _ := usecases.VoucherCodeUseCase.GetVoucherCodes(tempEcho.Context, payload)
					voucher.PromoCode = &voucherCodes[0].PromoCode
					expectedValue = models.RoundDown(*voucher.Validators.Value, 0)
				})

				Context("inquiry succeeded", func() {
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

				Context("inquiry error", func() {
					Context("voucher codes has been used", func() {
						var usedCif string

						JustBeforeEach(func() {
							response = models.Response{}
							reqpl.CIF = usedCif

							pl = reqpl
							e = test.NewDummyEcho(http.MethodPost, "/", pl)
							_ = handler.MultiRewardInquiry(e.Context)
							_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
							newExpectedResp(reqpl.CIF, true)
						})

						Context("with the same cif", func() {
							BeforeEach(func() {
								usedCif = reqpl.CIF
							})

							JustBeforeEach(func() {
								expectResp.Message = "Kode voucher sedang diproses"
							})

							It("expect response to return as expected", func() {
								Expect(response).To(Equal(expectResp))
							})
						})

						Context("with any other cif", func() {
							BeforeEach(func() {
								usedCif = "1122334455"
							})

							JustBeforeEach(func() {
								expectResp.Message = "Maaf, tidak ada reward yang tersedia"
							})

							It("expect response to return as expected", func() {
								Expect(response).To(Equal(expectResp))
							})
						})
					})

					Context("wrong voucher code", func() {
						BeforeEach(func() {
							usedCode := "something"
							voucher.PromoCode = &usedCode
						})

						JustBeforeEach(func() {
							newExpectedResp(reqpl.CIF, true)
							expectResp.Message = "Maaf, tidak ada reward yang tersedia"
						})

						It("expect response to return as expected", func() {
							Expect(response).To(Equal(expectResp))
						})
					})
				})
			})

			Context("referral", func() {
				var reward1 models.Reward
				var oslActive bool
				var withOslActive bool

				BeforeEach(func() {
					withEcho, oslActive, withOslActive = false, true, true
					campaign = fakedata.CampaignReferral()
					reward = fakedata.RewardDirectDisc(false)
					reward1 = fakedata.RewardIncentive(false)
					reward1.Validators.Channel = reward.Validators.Channel
					reward1.Validators.TransactionType = reward.Validators.TransactionType
					reward1.Validators.Product = reward.Validators.Product

					createReferral = models.RequestCreateReferral{
						CIF:    gofakeit.Regex("[1234567890]{10}"),
						Prefix: campaign.Metadata.Prefix,
					}
				})

				JustBeforeEach(func() {
					rewards = append(rewards, reward1, reward)
					campaign.Rewards = &rewards
					refCode = createCampaign(&campaign, createReferral)
					reqpl = mapReqpl(reqpl, reward, voucher, refCode)
					pl = reqpl
					e = test.NewDummyEcho(http.MethodPost, "/", pl)
					repos = config.NewRepositories(db.Sql, db.Bun)

					if withOslActive {
						mockCtrl := gomock.NewController(GinkgoT())
						mockRestRefRp := mocks.NewMockRestRefRepository(mockCtrl)
						restPl := models.ReqOslStatus{Cif: reqpl.CIF, ProductCode: reward.Validators.Product}
						mockRestRefRp.EXPECT().RRGetOslStatus(e.Context, restPl).Return(oslActive, nil)
						repos.ReferralsRestRepository = mockRestRefRp
					}

					usecases = config.NewUsecases(repos)
					handler.RewardUseCase = usecases.RewardUseCase
					_ = handler.MultiRewardInquiry(e.Context)
					_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
					newExpectedResp(reqpl.CIF)
				})

				Context("inquiry error osl active", func() {
					BeforeEach(func() {
						oslActive = false
					})

					JustBeforeEach(func() {
						newExpectedResp(reqpl.CIF, true)
						expectResp.Message = "OSL dari nasabah telah aktif"
					})

					It("expect response to return error", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("inquiry error max transaction", func() {
					BeforeEach(func() {
						withOslActive = false
						reward1.Validators.Incentive.MaxPerTransaction = helper.CreateFloat64(1)
					})

					JustBeforeEach(func() {
						newExpectedResp(reqpl.CIF, true)
						expectResp.Message = "Kode referral tidak dapat digunakan"
					})

					It("expect response to return error", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("inquiry succeeded", func() {
					JustBeforeEach(func() {
						newExpectedResp(reqpl.CIF)
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

				Context("inquiry error max per days / month", func() {
					BeforeEach(func() {
						withOslActive = false

						reward1.Validators.Incentive.MaxPerDay = helper.CreateFloat64(1)
						reward1.Validators.Incentive.MaxPerMonth = helper.CreateFloat64(1000000)
						reward1.Validators.Incentive.OslInactiveValidation = false
					})

					JustBeforeEach(func() {
						// successed inquiry
						responsData := response.Data.(map[string]interface{})
						plPayment.RootRefTrx = responsData["refTrx"].(string)
						plPayment.CIF = reqpl.CIF
						plPayment.RefCore = "11223344"
						_, _ = usecases.RewardUseCase.Payment(e.Context, &plPayment)

						// inquery
						g := test.NewDummyEcho(http.MethodPost, "/", pl)
						usecases = config.NewUsecases(repos)
						handler.RewardUseCase = usecases.RewardUseCase
						_ = handler.MultiRewardInquiry(g.Context)
						_ = json.Unmarshal(g.Response.Body.Bytes(), &response)
						newExpectedResp(reqpl.CIF, true)

						expectResp.Message = "Maaf, kode referral telah mencapai maksimal insentif"
					})

					It("expect response to return error", func() {
						Expect(response).To(Equal(expectResp))
					})
				})

				Context("error referral code used by same cif in same product", func() {
					BeforeEach(func() {
						withOslActive = false

						reward1.Validators.Incentive.OslInactiveValidation = false
					})

					JustBeforeEach(func() {
						responseData := response.Data.(map[string]interface{})
						plPayment.RootRefTrx = responseData["refTrx"].(string)
						plPayment.CIF = reqpl.CIF
						plPayment.RefCore = "1122334412"
						_, _ = usecases.RewardUseCase.Payment(e.Context, &plPayment)

						newExpectedResp(reqpl.CIF, true)
						usecases = config.NewUsecases(repos)
						handler.RewardUseCase = usecases.RewardUseCase
						f := test.NewDummyEcho(http.MethodPost, "/", pl)
						_ = handler.MultiRewardInquiry(f.Context)
						_ = json.Unmarshal(f.Response.Body.Bytes(), &response)
						expectResp.Message = "Kode referral tidak dapat digunakan pada produk yang sama"
					})

					It("expect response to return as expected", func() {
						Expect(response).To(Equal(expectResp))
					})
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

		Context("error min transaction amount", func() {
			BeforeEach(func() {
				reward = fakedata.RewardDirectDisc(true)
				reward.Validators.MinTransactionAmount = helper.CreateFloat64(100000)
				reqpl.TransactionAmount = helper.CreateFloat64(50000)
				expectedValue = models.RoundDown(*reward.Validators.Value, 0)
			})

			JustBeforeEach(func() {
				newExpectedResp(reqpl.CIF, true)
				expectResp.Message = "Maaf, tidak ada reward yang tersedia"
				expectResp.Description = "RewardIdx #0: transactionAmount untuk transaksi ini tidak valid untuk mendapatkan benefit"
			})

			It("expect to return expected error", func() {
				Expect(response).To(Equal(expectResp))
			})
		})

		Context("multiple errors happened", func() {
			BeforeEach(func() {
				reward = fakedata.RewardDirectDisc(true)
				reward1 := fakedata.RewardDirectDisc(true)
				reward.Validators.MinTransactionAmount = helper.CreateFloat64(100000)
				reward1.Validators.Channel = reward.Validators.Channel
				reward1.Validators.TransactionType = "something"
				reward1.Validators.Product = reward.Validators.Product
				reward1.PromoCode = reward.PromoCode

				rewards = append(rewards, reward1)

				reqpl.TransactionAmount = helper.CreateFloat64(50000)
				expectedValue = models.RoundDown(*reward.Validators.Value, 0)
			})

			JustBeforeEach(func() {
				newExpectedResp(reqpl.CIF, true)
				expectResp.Message = "Maaf, tidak ada reward yang tersedia"
				expectResp.Description = "RewardIdx #0: transactionType untuk transaksi ini tidak valid untuk mendapatkan benefit"
				expectResp.Description += ", RewardIdx #1: transactionAmount untuk transaksi ini tidak valid untuk mendapatkan benefit"
			})

			It("expect to return expected error", func() {
				Expect(response).To(Equal(expectResp))
			})
		})

		Context("inquiry an existing reward transaction", func() {
			BeforeEach(func() {
				reward = fakedata.RewardDirectDisc(true)
				expectedValue = models.RoundDown(*reward.Validators.Value, 0)
			})

			Context("using refTrx", func() {
				JustBeforeEach(func() {
					reqpl.RefTrx = responseData["refTrx"].(string)
					pl = reqpl
					response = models.Response{}
					e = test.NewDummyEcho(http.MethodPost, "/", pl)
					_ = handler.MultiRewardInquiry(e.Context)
					_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
					newExpectedResp(reqpl.CIF)
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("using same payload", func() {
				JustBeforeEach(func() {
					response = models.Response{}
					e = test.NewDummyEcho(http.MethodPost, "/", pl)
					_ = handler.MultiRewardInquiry(e.Context)
					_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
					newExpectedResp(reqpl.CIF)
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})

	Describe("GetRewardPromotions", func() {
		var reqpl models.RewardPromotionLists
		var rwdPromo map[string]interface{}

		BeforeEach(func() {
			reqpl = models.RewardPromotionLists{}
		})

		JustBeforeEach(func() {
			rewards = append(rewards, reward)
			campaign.Rewards = &rewards
			refCode = createCampaign(&campaign, createReferral)
			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.GetRewardPromotions(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
			rwdPromo = map[string]interface{}{}
			expectResp = models.Response{
				Code:    "00",
				Status:  "Success",
				Message: "Data Berhasil Dikirim",
				Data:    []interface{}{},
			}
		})

		Context("get list reward promotions", func() {

			Context("not empty", func() {
				BeforeEach(func() {
					reward = fakedata.RewardDirectDisc(true)
					reward.Validators.Product = "unique"
					reward.Validators.TransactionType = "unique"
					reward.Validators.Channel = "unique"

					reqpl = models.RewardPromotionLists{
						Product:         reward.Validators.Product,
						TransactionType: reward.Validators.TransactionType,
						Channel:         reward.Validators.Channel,
					}
				})

				JustBeforeEach(func() {
					response.Data.([]interface{})[0].(map[string]interface{})["id"] = 0

					rwdPromo = map[string]interface{}{
						"id":                 0,
						"name":               reward.Name,
						"description":        reward.Description,
						"termsAndConditions": reward.TermsAndConditions,
						"howToUse":           reward.HowToUse,
						"promoCode":          reward.PromoCode,
						"product":            reward.Validators.Product,
						"transactionType":    reward.Validators.TransactionType,
						"isPrivate":          "0",
					}
					expectResp = models.Response{
						Code:    "00",
						Status:  "Success",
						Message: "Data Berhasil Dikirim",
						Data:    []interface{}{rwdPromo},
					}
				})

				It("expect response to return as expected", func() {
					Expect(response.Data).To(ConsistOf(expectResp.Data))
				})
			})

			Context("empty", func() {
				BeforeEach(func() {
					reqpl = models.RewardPromotionLists{
						Product:         "unique-1",
						TransactionType: "unique-1",
						Channel:         "unique-1",
					}
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})
})

func inqResponse(responseData map[string]interface{}, reward models.Reward, value float64, cif string) map[string]interface{} {
	if reward == (models.Reward{}) || len(responseData) == 0 {
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
	if reward == (models.Reward{}) {
		return models.RespReferral{}
	}

	e := test.NewDummyEcho(http.MethodPost, "/", nil)
	repos := config.NewRepositories(db.Sql, db.Bun)
	usecases := config.NewUsecases(repos)

	_ = usecases.CampaignUseCase.CreateCampaign(e.Context, campaign)

	if createReferral == (models.RequestCreateReferral{}) {
		return models.RespReferral{}
	}

	rc, _ := usecases.ReferralsUseCase.UCreateReferralCodes(e.Context, createReferral)

	return rc
}

func newExpectedResp(cif string, isError ...bool) {
	if isError == nil {
		isError = append(isError, false)
	}

	expectResp = models.Response{
		Code:    "00",
		Status:  "Success",
		Message: "Data Berhasil Dikirim",
	}

	if isError[0] {
		response.Data = nil
		expectResp.Code = "99"
		expectResp.Status = "Error"
		expectResp.Message = ""

		return
	}

	if response.Data != nil {
		responseData = response.Data.(map[string]interface{})
	}

	expectResp.Data = inqResponse(responseData, reward, expectedValue, cif)
}
