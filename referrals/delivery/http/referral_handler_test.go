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

	"github.com/golang-migrate/migrate/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	db       *database.DbConfig
	migrator *migrate.Migrate

	_ = BeforeSuite(func() {
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
	var pl models.RequestCreateReferral
	var plGetHistory models.RequestHistoryIncentive
	var campaign models.Campaign
	var usecases config.Usecases

	config.LoadEnv()

	BeforeEach(func() {
		response, expectResp = models.Response{}, models.Response{}
		e = test.NewDummyEcho(http.MethodPost, "/", pl)
		_, usecases = test.LoadRealRepoUsecase(db)
		handler.ReferralUseCase = usecases.ReferralsUseCase
		pl = models.RequestCreateReferral{CIF: "1122334455", Prefix: "PSE"}
		plGetHistory = models.RequestHistoryIncentive{RefCif: "123456789"}
	})

	Describe("HGenerateReferralCodes", func() {
		JustBeforeEach(func() {
			e = test.NewDummyEcho(http.MethodPost, "/", pl)
			_ = handler.HGenerateReferralCodes(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		BeforeEach(func() {
			campaign = fakedata.Campaign()
			campaign.Status = helper.CreateInt8(1)
			campaign.Rewards = &[]models.Reward{fakedata.Reward(), fakedata.Reward()}
			_ = usecases.CampaignUseCase.CreateCampaign(e.Context, &campaign)

			pl.Prefix = campaign.Metadata.Prefix
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
					pl.CIF = "1122334466"
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
					pl.CIF = ""
					expectResp.Message = "Key: 'RequestCreateReferral.CIF' Error:Field validation for 'CIF' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("blank prefix", func() {
				BeforeEach(func() {
					pl.Prefix = ""
					expectResp.Message = "Key: 'RequestCreateReferral.Prefix' Error:Field validation for 'Prefix' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("prefix character greater than 5 character", func() {
				BeforeEach(func() {
					pl.Prefix = "QWEASD"
					expectResp.Message = "Key: 'RequestCreateReferral.Prefix' Error:Field validation for 'Prefix' failed on the 'max' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("CIF character greater than 10 character", func() {
				BeforeEach(func() {
					pl.CIF = "112233445566"
					expectResp.Message = "Key: 'RequestCreateReferral.CIF' Error:Field validation for 'CIF' failed on the 'max' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})

	Describe("HGetHistoriesIncentive", func() {
		JustBeforeEach(func() {
			e = test.NewDummyEcho(http.MethodPost, "/", plGetHistory)
			_ = handler.HGetHistoriesIncentive(e.Context)
			_ = json.Unmarshal(e.Response.Body.Bytes(), &response)
		})

		Context("get referral incentive error", func() {
			BeforeEach(func() {
				expectResp = models.Response{
					Code:   "99",
					Status: "Error",
				}
			})

			Context("when cif is empty get referral incentive history", func() {
				BeforeEach(func() {
					plGetHistory.RefCif = ""
					expectResp.Message = "Key: 'RequestHistoryIncentive.RefCif' Error:Field validation for 'RefCif' failed on the 'required' tag"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})

			Context("when cif has no referral incentive history", func() {
				BeforeEach(func() {
					plGetHistory.RefCif = "1017441370"
					expectResp.Message = "Data history incentive referral tidak ditemukan"
				})

				It("expect response to return as expected", func() {
					Expect(response).To(Equal(expectResp))
				})
			})
		})
	})
})
