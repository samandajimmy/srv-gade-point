package config

import (
	_campaignUseCase "gade/srv-gade-point/campaigns/usecase"
	_campaignTrxUseCase "gade/srv-gade-point/campaigntrxs/usecase"
	_pHistoryUseCase "gade/srv-gade-point/pointhistories/usecase"
	_quotaUseCase "gade/srv-gade-point/quotas/usecase"
	_referralsUseCase "gade/srv-gade-point/referrals/usecase"
	_referralTrxUseCase "gade/srv-gade-point/referraltrxs/usecase"
	_rewardUseCase "gade/srv-gade-point/rewards/usecase"
	_rewardTrxUC "gade/srv-gade-point/rewardtrxs/usecase"
	_tagUseCase "gade/srv-gade-point/tags/usecase"
	_tokenUseCase "gade/srv-gade-point/tokens/usecase"
	_userUseCase "gade/srv-gade-point/users/usecase"
	_voucherCodeUseCase "gade/srv-gade-point/vouchercodes/usecase"
	_voucherUseCase "gade/srv-gade-point/vouchers/usecase"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/campaigntrxs"
	"gade/srv-gade-point/pointhistories"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/referrals"
	"gade/srv-gade-point/referraltrxs"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/rewardtrxs"
	"gade/srv-gade-point/tags"
	"gade/srv-gade-point/tokens"
	"gade/srv-gade-point/users"
	"gade/srv-gade-point/vouchercodes"
	"gade/srv-gade-point/vouchers"
)

type Usecases struct {
	CampaignUseCase    campaigns.UseCase
	CampaignTrxUseCase campaigntrxs.UseCase
	PHistoryUseCase    pointhistories.UseCase
	QuotaUseCase       quotas.QUseCase
	ReferralsUseCase   referrals.RefUseCase
	ReferralTrxUseCase referraltrxs.UseCase
	RewardUseCase      rewards.UseCase
	RewardTrxUC        rewardtrxs.UseCase
	TagUseCase         tags.TUseCase
	TokenUseCase       tokens.UseCase
	UserUseCase        users.UseCase
	VoucherCodeUseCase vouchercodes.UseCase
	VoucherUseCase     vouchers.VUsecase
}

func NewUsecases(repos Repositories) Usecases {
	tokenUseCase := _tokenUseCase.NewTokenUseCase(repos.TokenRepository)
	tagUseCase := _tagUseCase.NewTagUseCase(repos.TagRepository)
	pHistoryUseCase := _pHistoryUseCase.NewPointHistoryUseCase(repos.PHistoryRepository)
	rewardTrxUseCase := _rewardTrxUC.NewRewardtrxUseCase(repos.RewardTrxRepository)
	referralTrxUseCase := _referralTrxUseCase.NewReferralTrxUseCase(repos.ReferralTrxRepository)
	quotaUseCase := _quotaUseCase.NewQuotaUseCase(repos.QuotaRepository, rewardTrxUseCase)
	voucherCodeUseCase := _voucherCodeUseCase.NewVoucherCodeUseCase(repos.VoucherCodeRepository,
		repos.VoucherRepository)
	voucherUseCase := _voucherUseCase.NewVoucherUseCase(repos.VoucherRepository,
		repos.CampaignRepository, repos.PHistoryRepository, repos.VoucherCodeRepository)
	campaignTrxUseCase := _campaignTrxUseCase.NewCampaignTrxUseCase(repos.CampaignTrxRepository)
	userUseCase := _userUseCase.NewUserUseCase(repos.UserRepository)
	referralsUseCase := _referralsUseCase.NewReferralUseCase(repos.ReferralsRepository, repos.CampaignRepository)
	rewardUseCase := _rewardUseCase.NewRewardUseCase(repos.RewardRepository, repos.CampaignRepository,
		tagUseCase, quotaUseCase, voucherUseCase, repos.VoucherCodeRepository, repos.RewardTrxRepository,
		repos.ReferralTrxRepository, referralsUseCase)
	campaignUseCase := _campaignUseCase.NewCampaignUseCase(repos.CampaignRepository, rewardUseCase)

	return Usecases{
		CampaignUseCase:    campaignUseCase,
		CampaignTrxUseCase: campaignTrxUseCase,
		PHistoryUseCase:    pHistoryUseCase,
		QuotaUseCase:       quotaUseCase,
		ReferralsUseCase:   referralsUseCase,
		ReferralTrxUseCase: referralTrxUseCase,
		RewardUseCase:      rewardUseCase,
		RewardTrxUC:        rewardTrxUseCase,
		TagUseCase:         tagUseCase,
		TokenUseCase:       tokenUseCase,
		UserUseCase:        userUseCase,
		VoucherCodeUseCase: voucherCodeUseCase,
		VoucherUseCase:     voucherUseCase,
	}
}
