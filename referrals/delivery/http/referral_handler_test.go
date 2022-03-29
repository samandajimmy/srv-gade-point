package http_test

import (
	"encoding/json"
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	refhttp "gade/srv-gade-point/referrals/delivery/http"
	"gade/srv-gade-point/test"
	"gade/srv-gade-point/test/fakedata"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var (
	db       *database.DbConfig
	migrator *migrate.Migrate

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

var _ = Describe("ReferralHandler", func() {
	var handler refhttp.ReferralHandler
	var e test.DummyEcho
	var response models.Response
	var expectResp models.Response
	var pl interface{}
	var campaign models.Campaign
	var usecases config.Usecases

	BeforeEach(func() {
		response, expectResp = models.Response{}, models.Response{}
		e = test.NewDummyEcho(http.MethodPost, "/", pl)
		_, usecases = test.LoadRealRepoUsecase(db)
		handler.ReferralUseCase = usecases.ReferralsUseCase
	})

	Describe("HGenerateReferralCodes", func() {
		var reqpl models.RequestCreateReferral

		JustBeforeEach(func() {
			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.HGenerateReferralCodes(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		BeforeEach(func() {
			withPromoCode := true
			campaign := fakedata.CampaignReferral()
			campaign.Rewards = &[]models.Reward{fakedata.Reward(withPromoCode), fakedata.Reward(withPromoCode)}
			_ = usecases.CampaignUseCase.CreateCampaign(e.Context, &campaign)
			reqpl = models.RequestCreateReferral{CIF: "1122334455", Prefix: campaign.Metadata.Prefix}
		})

		Context("create referral code succeeded", func() {
			JustBeforeEach(func() {
				data := response.Data.(map[string]interface{})
				data["referralCode"] = ""
				response.Data = data
			})

			BeforeEach(func() {
				expectResp = models.Response{
					Code:    "00",
					Status:  "Success",
					Message: "Data Berhasil Dikirim",
					Data: map[string]interface{}{
						"cif":          "1122334455",
						"referralCode": "",
					},
				}
			})

			It("expect response to return as expected", func() {
				Expect(response).To(Equal(expectResp))
			})
		})

		Context("create referral code error", func() {
			BeforeEach(func() {
				expectResp = models.Response{
					Code:   "99",
					Status: "Error",
				}
			})

			Context("duplicated CIF", func() {
				BeforeEach(func() {
					reqpl.CIF = "1122334466"
					expectResp.Message = "Cif 1122334466 sudah memiliki Kode Referral"
				})

				JustBeforeEach(func() {
					response = models.Response{}
					e = test.NewDummyEcho(http.MethodPost, "/", pl)
					_ = handler.HGenerateReferralCodes(e.Context)
					_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("blank CIF", func() {
				BeforeEach(func() {
					reqpl.CIF = ""
					expectResp.Message = "Key: 'RequestCreateReferral.CIF' Error:Field validation for 'CIF' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("blank prefix", func() {
				BeforeEach(func() {
					reqpl.Prefix = ""
					expectResp.Message = "Key: 'RequestCreateReferral.Prefix' Error:Field validation for 'Prefix' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("prefix character greater than 5 character", func() {
				BeforeEach(func() {
					reqpl.Prefix = "QWEASD"
					expectResp.Message = "Key: 'RequestCreateReferral.Prefix' Error:Field validation for 'Prefix' failed on the 'max' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("CIF character greater than 10 character", func() {
				BeforeEach(func() {
					reqpl.CIF = "112233445566"
					expectResp.Message = "Key: 'RequestCreateReferral.CIF' Error:Field validation for 'CIF' failed on the 'max' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})

	Describe("HGetReferralCodes", func() {
		var reqpl models.RequestReferralCodeUser
		var usedCif = gofakeit.Regex("[1234567890]{10}")
		var refCode models.RespReferral
		var plInquiry models.PayloadValidator
		var plPayment models.RewardPayment

		JustBeforeEach(func() {
			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.HGetReferralCodes(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		BeforeEach(func() {
			_ = viper.UnmarshalKey("reward.inquiry", &plInquiry)
			_ = viper.UnmarshalKey("reward.payment", &plPayment)
			withPromoCode := false
			// create voucher
			voucher := fakedata.VoucherDirectDisc()
			_ = usecases.VoucherUseCase.CreateVoucher(e.Context, &voucher)
			// create campaign referral
			campaign = fakedata.CampaignReferral()
			reward1 := fakedata.RewardDirectDisc(withPromoCode)
			reward2 := fakedata.RewardIncentive(withPromoCode)
			campaign.Rewards = &[]models.Reward{reward1, reward2}
			_ = usecases.CampaignUseCase.CreateCampaign(e.Context, &campaign)
			// create referral code
			createReferral := models.RequestCreateReferral{
				CIF:    usedCif,
				Prefix: campaign.Metadata.Prefix,
			}
			refCode, _ = usecases.ReferralsUseCase.UCreateReferralCodes(e.Context, createReferral)
			// prepare reward inquiry data
			plInquiry.CIF = usedCif
			plInquiry.IsMulti = true
			plInquiry.TransactionDate = time.Now().AddDate(0, 0, 0).Format(models.DateTimeFormat + ".000")
			plInquiry.PromoCode = refCode.ReferralCode
			plInquiry.Validators.Channel = reward1.Validators.Channel
			plInquiry.Validators.Product = reward1.Validators.Product
			plInquiry.Validators.TransactionType = reward1.Validators.TransactionType
			// get the reward
			inquiry, _ := usecases.RewardUseCase.Inquiry(e.Context, &plInquiry)
			plPayment.RootRefTrx = inquiry.RefTrx
			_, _ = usecases.RewardUseCase.Payment(e.Context, &plPayment)

			reqpl = models.RequestReferralCodeUser{CIF: usedCif}
		})

		Context("get referral code detail succeeded", func() {
			BeforeEach(func() {
				expectResp = models.Response{
					Code:    "00",
					Status:  "Success",
					Message: "Data Berhasil Dikirim",
					Data: map[string]interface{}{
						"referralCode": refCode.ReferralCode,
						"incentive": map[string]interface{}{
							"isValid":       true,
							"perDay":        float64(155000),
							"perMonth":      float64(155000),
							"total":         float64(0),
							"validPerDay":   true,
							"validPerMonth": true,
						},
					},
				}
			})

			It("expect response to return as expected", func() {
				Expect(response).To(Equal(expectResp))
			})
		})
	})

	Describe("HGetHistoriesIncentive", func() {
		var reqpl models.RequestHistoryIncentive
		var refCif = gofakeit.Regex("[1234567890]{10}")
		var userCif = gofakeit.Regex("[1234567890]{10}")
		var refCode models.RespReferral
		var plInquiry models.PayloadValidator
		var plPayment models.RewardPayment

		JustBeforeEach(func() {
			pl = reqpl
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.HGetHistoriesIncentive(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		BeforeEach(func() {
			_ = viper.UnmarshalKey("reward.inquiry", &plInquiry)
			_ = viper.UnmarshalKey("reward.payment", &plPayment)
			withPromoCode := false

			// create voucher
			voucher := fakedata.VoucherDirectDisc()
			_ = usecases.VoucherUseCase.CreateVoucher(e.Context, &voucher)

			// create campaign referral
			campaign = fakedata.CampaignReferral()
			reward1 := fakedata.RewardDirectDisc(withPromoCode)
			reward2 := fakedata.RewardIncentive(withPromoCode)
			campaign.Rewards = &[]models.Reward{reward1, reward2}
			_ = usecases.CampaignUseCase.CreateCampaign(e.Context, &campaign)

			// create referral code
			createReferral := models.RequestCreateReferral{
				CIF:    refCif,
				Prefix: campaign.Metadata.Prefix,
			}
			refCode, _ = usecases.ReferralsUseCase.UCreateReferralCodes(e.Context, createReferral)

			// prepare reward inquiry data
			plInquiry.CIF = userCif
			plInquiry.IsMulti = true
			plInquiry.TransactionDate = time.Now().AddDate(0, 0, 0).Format(models.DateTimeFormat + ".000")
			plInquiry.PromoCode = refCode.ReferralCode
			plInquiry.Validators.Channel = reward1.Validators.Channel
			plInquiry.Validators.Product = reward1.Validators.Product
			plInquiry.Validators.TransactionType = reward1.Validators.TransactionType
			plInquiry.TransactionAmount = helper.CreateFloat64(100000)

			// get the reward
			inquiry, _ := usecases.RewardUseCase.Inquiry(e.Context, &plInquiry)

			plPayment.RootRefTrx = inquiry.RefTrx
			_, _ = usecases.RewardUseCase.Payment(e.Context, &plPayment)

			reqpl.Limit = 1
			reqpl.RefCif = refCif
		})

		Context("get referral incentive history succeed", func() {
			BeforeEach(func() {
				expectResp.Code = "00"
				expectResp.Status = "Success"
			})

			Context("when payload page is empty get referral incentive history", func() {
				It("expect response to return as expected", func() {
					Expect(response.Code).To(Equal(expectResp.Code))
					Expect(response.Status).To(Equal(expectResp.Status))
					Expect(response.Data).ToNot(BeNil())
				})
			})

			Context("when payload page is not empty get referral incentive history", func() {
				BeforeEach(func() {
					reqpl.Page = 1
				})

				It("expect response to return as expected", func() {
					Expect(response.Code).To(Equal(expectResp.Code))
					Expect(response.Status).To(Equal(expectResp.Status))
					Expect(response.Data).ToNot(BeNil())
				})
			})
		})

		Context("get referral incentive error", func() {
			BeforeEach(func() {
				expectResp = models.Response{
					Code:   "99",
					Status: "Error",
				}
			})

			Context("when limit is empty referral incentive history", func() {
				BeforeEach(func() {
					reqpl.RefCif = "1017441370"
					reqpl.Limit = 0
					expectResp.Message = "Key: 'RequestHistoryIncentive.Limit' Error:Field validation for 'Limit' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("when cif is empty get referral incentive history", func() {
				BeforeEach(func() {
					reqpl.RefCif = ""
					reqpl.Limit = 1
					expectResp.Message = "Key: 'RequestHistoryIncentive.RefCif' Error:Field validation for 'RefCif' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("when cif has no referral incentive history", func() {
				BeforeEach(func() {
					reqpl.RefCif = "1017441370"
					reqpl.Limit = 1
					expectResp.Message = "Data history incentive referral tidak ditemukan"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})
})
