package usecase_test

import (
	"gade/srv-gade-point/config"
	. "gade/srv-gade-point/helper"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/rewards/usecase"
	"gade/srv-gade-point/test"
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("RewardInqUcase", func() {
	var rUS rewards.UseCase
	var e test.DummyEcho
	var mockCtrl *gomock.Controller
	var mockUsecases mocks.MockUsecases
	var mockRepos mocks.MockRepositories
	var mockCampaigns []*models.Campaign
	var mockRespInquiry models.RewardsInquiry
	var mockSumIncentive models.SumIncentive

	var pl = models.PayloadValidator{
		IsMulti:           true,
		CIF:               "11223344515111",
		Referrer:          "1122334451522",
		Phone:             "081511150290",
		PromoCode:         "cacing",
		CustomerName:      "Jimmy Samanda Rasu",
		TransactionDate:   "2022-01-18 16:19:32.121",
		TransactionAmount: CreateFloat64(100000),
		Validators: &models.Validator{
			Channel:         "9997",
			Product:         "37",
			TransactionType: "OP",
		},
		CodeType: "referral",
	}

	config.LoadEnv()
	config.LoadTestData()

	BeforeEach(func() {
		e = test.NewDummyEcho(http.MethodPost, "/")
		mockCtrl = gomock.NewController(GinkgoT())
		mockUsecases = mocks.NewMockUsecases(mockCtrl)
		mockRepos = mocks.NewMockRepository(mockCtrl)
	})

	Describe("Inquiry", func() {
		var campaign models.Campaign
		var rewards []models.Reward
		var respData models.RewardsInquiry
		var errs *models.ResponseErrors

		BeforeEach(func() {
			_ = viper.UnmarshalKey("campaign.referralcampaign", &mockCampaigns)
			_ = viper.UnmarshalKey("reward.responseinquiry", &mockRespInquiry)
			_ = viper.UnmarshalKey("incentive.sumincentive", &mockSumIncentive)

			campaign = *mockCampaigns[0]
			rewards = *campaign.Rewards
			rewards[0].Campaign = &campaign
			rewards[1].Campaign = &campaign
		})

		JustBeforeEach(func() {
			rUS = usecase.NewRewardUseCase(mockRepos.MockRRp, mockRepos.MockCRp, mockUsecases.MockTUs,
				mockUsecases.MockQUs, mockUsecases.MockVUs, mockRepos.MockVcRp, mockRepos.MockRtRp,
				mockRepos.MockRefTRp, mockUsecases.MockRefUs)

			respData, errs = rUS.Inquiry(e.Context, &pl)
		})

		Context("promo referral code", func() {
			Context("succeeded", func() {
				BeforeEach(func() {
					mockUsecases.MockVUs.EXPECT().GetVoucherCode(e.Context, &pl, false).Return(nil, "", nil)
					mockRepos.MockCRp.EXPECT().GetReferralCampaign(e.Context, pl).Return(&mockCampaigns)
					mockRepos.MockRtRp.EXPECT().GetRewardByPayload(e.Context, pl, nil).Return([]*models.Reward{}, nil)
					mockUsecases.MockRefUs.EXPECT().UValidateReferrer(e.Context, pl, &campaign).Return(mockSumIncentive, nil)
					mockRepos.MockRRp.EXPECT().GetRewardTags(e.Context, &rewards[0]).Return(nil, nil)
					mockRepos.MockRRp.EXPECT().GetRewardTags(e.Context, &rewards[1]).Return(nil, nil)
					mockUsecases.MockQUs.EXPECT().CheckQuota(e.Context, rewards[0], &pl).Return(true, nil)
					mockUsecases.MockQUs.EXPECT().CheckQuota(e.Context, rewards[1], &pl).Return(true, nil)
					mockUsecases.MockQUs.EXPECT().UpdateReduceQuota(e.Context, rewards[0].ID).Return(nil)
					mockUsecases.MockQUs.EXPECT().UpdateReduceQuota(e.Context, rewards[1].ID).Return(nil)
					mockRepos.MockRRp.EXPECT().RGetRandomId(20).Return("someIDs").Times(3)
					mockRepos.MockRtRp.EXPECT().Create(e.Context, pl, mockRespInquiry).Return(nil, nil)
				})

				It("expect to return reward inquiry response", func() {
					Expect(respData).To(Equal(mockRespInquiry))
				})

				It("expect to return nil response errors", func() {
					Expect(errs).To(Equal(&models.ResponseErrors{}))
				})
			})

			Context("exceeded incentive", func() {
				BeforeEach(func() {
					mockSumIncentive.IsValid = false

					mockUsecases.MockVUs.EXPECT().GetVoucherCode(e.Context, &pl, false).Return(nil, "", nil)
					mockRepos.MockCRp.EXPECT().GetReferralCampaign(e.Context, pl).Return(&mockCampaigns)
					mockRepos.MockRtRp.EXPECT().GetRewardByPayload(e.Context, pl, nil).Return([]*models.Reward{}, nil)
					mockUsecases.MockRefUs.EXPECT().UValidateReferrer(e.Context, pl, &campaign).Return(mockSumIncentive, nil)
				})

				It("expect to return string `Kode referral telah mencapai batas maksimal`", func() {
					Expect(errs.Title).To(Equal(models.ErrRefTrxExceeded.Error()))
				})

				It("expect to return response data nil", func() {
					Expect(respData).To(Equal(models.RewardsInquiry{}))
				})
			})
		})
	})
})
