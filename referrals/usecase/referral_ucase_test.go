package usecase_test

import (
	"net/http"

	"gade/srv-gade-point/config"
	. "gade/srv-gade-point/helper"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"gade/srv-gade-point/referrals/usecase"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("ReferralUcase", func() {
	var refUS referrals.RefUseCase
	var e DummyEcho
	var mockCtrl *gomock.Controller
	var mockRepos mocks.MockRepositories

	config.LoadEnv()
	config.LoadTestData()

	BeforeEach(func() {
		e = NewDummyEcho(http.MethodPost, "/")
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepos = mocks.NewMockRepository(mockCtrl)
	})

	Describe("UValidateReferrer", func() {
		It("happy case", func() {
			var mockCampaigns []models.Campaign
			var mockSumIncentive models.SumIncentive
			var pl = models.PayloadValidator{
				PromoCode:       "cacing",
				TransactionDate: "2022-01-18 16:19:32.121",
			}

			_ = viper.UnmarshalKey("campaign.referralcampaign", &mockCampaigns)
			_ = viper.UnmarshalKey("incentive.sumincentive", &mockSumIncentive)

			rewards := *mockCampaigns[0].Rewards
			reward := rewards[1]
			mockRepos.MockRefRp.EXPECT().RSumRefIncentive(e.Context, pl.PromoCode, reward).Return(mockSumIncentive, nil)

			refUS = usecase.NewReferralUseCase(mockRepos.MockRefRp, mockRepos.MockCRp)
			data, err := refUS.UValidateReferrer(e.Context, pl, &mockCampaigns[0])

			Expect(ToJson(data)).To(Equal(ToJson(mockSumIncentive)))
			Expect(err).To(BeNil())
		})
	})
})
