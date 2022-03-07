package config

import (
	_campaignRepository "gade/srv-gade-point/campaigns/repository"
	_campaignTrxRepository "gade/srv-gade-point/campaigntrxs/repository"
	_pHistoryRepository "gade/srv-gade-point/pointhistories/repository"
	_quotaRepository "gade/srv-gade-point/quotas/repository"
	_referralsRepository "gade/srv-gade-point/referrals/repository"
	_referralTrxRepository "gade/srv-gade-point/referraltrxs/repository"
	_rewardRepository "gade/srv-gade-point/rewards/repository"
	_rewardTrxRepository "gade/srv-gade-point/rewardtrxs/repository"
	_tagRepository "gade/srv-gade-point/tags/repository"
	_tokenRepository "gade/srv-gade-point/tokens/repository"
	_userRepository "gade/srv-gade-point/users/repository"
	_voucherCodeRepository "gade/srv-gade-point/vouchercodes/repository"
	_voucherRepository "gade/srv-gade-point/vouchers/repository"

	"database/sql"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/campaigntrxs"
	"gade/srv-gade-point/database"
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

type Repositories struct {
	CampaignRepository    campaigns.CRepository
	CampaignTrxRepository campaigntrxs.Repository
	PHistoryRepository    pointhistories.Repository
	QuotaRepository       quotas.Repository
	ReferralsRepository   referrals.RefRepository
	ReferralTrxRepository referraltrxs.RefTRepository
	RewardRepository      rewards.RRepository
	RewardTrxRepository   rewardtrxs.RtRepository
	TagRepository         tags.Repository
	TokenRepository       tokens.Repository
	UserRepository        users.Repository
	VoucherCodeRepository vouchercodes.VcRepository
	VoucherRepository     vouchers.Repository
}

func NewRepositories(dbConn *sql.DB, dbBun *database.DbBun) Repositories {
	tokenRepository := _tokenRepository.NewPsqlTokenRepository(dbConn)
	tagRepository := _tagRepository.NewPsqlTagRepository(dbConn)
	pHistoryRepository := _pHistoryRepository.NewPsqlPointHistoryRepository(dbConn)
	rewardTrxRepository := _rewardTrxRepository.NewPsqlRewardTrxRepository(dbConn)
	referralTrxRepository := _referralTrxRepository.NewPsqlReferralTrxRepository(dbConn)
	quotaRepository := _quotaRepository.NewPsqlQuotaRepository(dbConn)
	voucherRepository := _voucherRepository.NewPsqlVoucherRepository(dbConn, dbBun)
	voucherCodeRepository := _voucherCodeRepository.NewPsqlVoucherCodeRepository(dbConn)
	rewardRepository := _rewardRepository.NewPsqlRewardRepository(dbConn, dbBun, quotaRepository, tagRepository)
	campaignRepository := _campaignRepository.NewPsqlCampaignRepository(dbConn, dbBun, rewardRepository)
	campaignTrxRepository := _campaignTrxRepository.NewPsqlCampaignTrxRepository(dbConn)
	userRepository := _userRepository.NewPsqlUserRepository(dbConn)
	referralsRepository := _referralsRepository.NewPsqlReferralRepository(dbConn, dbBun)

	return Repositories{
		CampaignRepository:    campaignRepository,
		CampaignTrxRepository: campaignTrxRepository,
		PHistoryRepository:    pHistoryRepository,
		QuotaRepository:       quotaRepository,
		ReferralsRepository:   referralsRepository,
		ReferralTrxRepository: referralTrxRepository,
		RewardRepository:      rewardRepository,
		RewardTrxRepository:   rewardTrxRepository,
		TagRepository:         tagRepository,
		TokenRepository:       tokenRepository,
		UserRepository:        userRepository,
		VoucherCodeRepository: voucherCodeRepository,
		VoucherRepository:     voucherRepository,
	}
}
