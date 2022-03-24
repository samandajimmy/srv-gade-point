package usecase_test

import (
	"net/http"

	. "gade/srv-gade-point/helper"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"gade/srv-gade-point/referrals/usecase"
	"gade/srv-gade-point/test"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("ReferralUcase", func() {
	var refUS referrals.RefUseCase
	var e test.DummyEcho
	var mockCtrl *gomock.Controller
	var mockRepos mocks.MockRepositories
	var mockRefCode models.ReferralCodes
	var err error

	BeforeEach(func() {
		e = test.NewDummyEcho(http.MethodPost, "/")
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepos, _ = e.LoadMockRepoUsecase(mockCtrl)
		refUS = usecase.NewReferralUseCase(mockRepos.MockRefRp, mockRepos.MockCRp)

		_ = viper.UnmarshalKey("referralcode.one", &mockRefCode)
	})

	Describe("UCreateReferralCodes", func() {
		var output models.RespReferral
		var expetedOutput models.RespReferral
		var pl = models.RequestCreateReferral{CIF: "0987654325", Prefix: "PSE"}

		JustBeforeEach(func() {
			output, err = refUS.UCreateReferralCodes(e.Context, pl)
		})

		Context("when the cif already exist", func() {
			BeforeEach(func() {
				refCode := models.ReferralCodes{
					CIF:        pl.CIF,
					CampaignId: int64(1),
				}
				code := "PSE11223"

				mockRepos.MockRefRp.EXPECT().RGetReferralByCif(e.Context, pl.CIF).Return(models.ReferralCodes{}, nil)
				mockRepos.MockRefRp.EXPECT().RGetCampaignId(e.Context, pl.Prefix).Return(int64(1), nil)
				mockRepos.MockRefRp.EXPECT().RGenerateCode(e.Context, refCode, pl.Prefix).Return(code)
				refCode.ReferralCode = code
				mockRepos.MockRefRp.EXPECT().RCreateReferral(e.Context, refCode).Return(mockRefCode, nil)
				expetedOutput = models.RespReferral{
					CIF:          mockRefCode.CIF,
					ReferralCode: mockRefCode.ReferralCode,
				}
			})

			It("expect to return nil err", func() {
				Expect(err).To(BeNil())
			})

			It("expect to return referral code data", func() {
				Expect(output).To(Equal(expetedOutput))
			})
		})
	})

	Describe("UValidateReferrer", func() {
		It("happy case", func() {
			var mockCampaign models.Campaign
			var mockSumIncentive models.SumIncentive
			var pl = models.PayloadValidator{
				PromoCode:       "cacing",
				TransactionDate: "2022-01-18 16:19:32.121",
			}

			_ = viper.UnmarshalKey("campaign.referralcampaign", &mockCampaign)
			_ = viper.UnmarshalKey("incentive.sumincentive", &mockSumIncentive)

			rewards := *mockCampaign.Rewards
			reward := rewards[1]
			mockCampaignReferral := models.CampaignReferral{
				Campaign:    mockCampaign,
				CifReferrer: "",
			}
			mockRepos.MockRefRp.EXPECT().RSumRefIncentive(e.Context, pl.PromoCode, reward).Return(mockSumIncentive, nil)

			refUS = usecase.NewReferralUseCase(mockRepos.MockRefRp, mockRepos.MockCRp)
			data, err := refUS.UValidateReferrer(e.Context, pl, &mockCampaignReferral)

			Expect(ToJson(data)).To(Equal(ToJson(mockSumIncentive)))
			Expect(err).To(BeNil())
		})
	})
})
