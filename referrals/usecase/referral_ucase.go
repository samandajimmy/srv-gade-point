package usecase

import (
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"time"

	"github.com/labstack/echo"
)

type referralUseCase struct {
	referralRepo referrals.RefRepository
	campaignRepo campaigns.CRepository
	restReferral referrals.RestRefRepository
}

// NewReferralUseCase will create new an rewardtrxUseCase object representation of referrals.RefUseCase interface
func NewReferralUseCase(
	refRepo referrals.RefRepository,
	campaignRepo campaigns.CRepository,
	restReferral referrals.RestRefRepository,
) referrals.RefUseCase {
	return &referralUseCase{
		referralRepo: refRepo,
		campaignRepo: campaignRepo,
		restReferral: restReferral,
	}
}

func (rcUc *referralUseCase) UCreateReferralCodes(c echo.Context, requestReferralCodes models.RequestCreateReferral) (models.RespReferral, error) {
	var payload models.ReferralCodes

	// Check if a cif already had a referral code
	referral, err := rcUc.referralRepo.RGetReferralByCif(c, requestReferralCodes.CIF)

	if err != nil {
		return models.RespReferral{}, err
	}

	// if exist, return existing referral code
	if referral.ReferralCode != "" {
		return models.RespReferral{}, models.DynamicErr(models.CifRefCodeExisted, requestReferralCodes.CIF)
	}

	// if not, generate referral code, getting campaign id first
	campaignId, err := rcUc.getCampaignIdByPrefix(c, requestReferralCodes.Prefix)

	if err != nil {
		logger.Make(c).Debug(err)

		return models.RespReferral{}, err
	}

	// map campaign id, generate referral code according to requirements (with prefix)
	payload.CampaignId = campaignId
	payload.CIF = requestReferralCodes.CIF
	payload.ReferralCode = rcUc.referralRepo.RGenerateCode(c, payload, requestReferralCodes.Prefix)
	// create referral code data
	result, err := rcUc.referralRepo.RCreateReferral(c, payload)

	if err != nil {
		return models.RespReferral{}, models.ErrGenerateRefCodeFailed
	}

	return models.RespReferral{
		CIF:          result.CIF,
		ReferralCode: result.ReferralCode,
	}, nil
}

func (rcUc *referralUseCase) UGetReferralCodes(c echo.Context, requestGetReferral models.RequestReferralCodeUser) (models.RespReferralDetail, error) {
	var result models.RespReferralDetail

	// get referral codes data by cif
	referral, err := rcUc.referralRepo.RGetReferralByCif(c, requestGetReferral.CIF)

	// throw empty data if referral codes not found
	if err != nil {
		logger.Make(c).Debug(models.ErrRefCodesNF)

		return result, models.ErrRefCodesNF
	}

	if referral.ReferralCode == "" {
		return result, models.ErrRefCodesNF
	}

	// get reward incentive for referrer
	reward, err := rcUc.campaignRepo.GetRewardIncentiveByCampaign(c, referral.CampaignId)

	if referral.ReferralCode == "" {
		return result, err
	}

	// sum the incentive based on the ref trx (per day and per month)
	objIncentive, err := rcUc.referralRepo.RSumRefIncentive(c, referral.ReferralCode)

	if err != nil {
		return result, err
	}

	// validate the max per day and per month
	validator := reward.Validators
	validator.Incentive.Validate(&objIncentive)

	// if data referral codes found send result
	result.ReferralCode = referral.ReferralCode
	result.Incentive = objIncentive.DetailIncentive

	return result, nil
}

func (ref *referralUseCase) UValidateReferrer(c echo.Context, pl models.PayloadValidator, campaignReferral *models.CampaignReferral) (models.ObjIncentive, error) {
	var objIncentive models.ObjIncentive
	var chosenReward models.Reward
	var err error

	// get ref trx based on the active ref campaign
	// and reward data with used ref code
	if campaignReferral == nil {
		campaignReferral = ref.campaignRepo.GetReferralCampaign(c, pl)
	}

	for _, reward := range *campaignReferral.Rewards {
		if reward.Validators.Target != models.RefTargetReferrer {
			continue
		}

		// sum the incentive based on the ref trx (per day and per month)
		objIncentive, err = ref.referralRepo.RSumRefIncentive(c, pl.PromoCode)

		if err != nil {
			continue
		}

		chosenReward = reward

	}

	// validate the max per day and per month
	validator := chosenReward.Validators

	if validator == nil {
		return objIncentive, nil
	}

	objIncentive.PerTransaction = *pl.TransactionAmount
	validator.Incentive.Validate(&objIncentive)

	// when its not valid then no need to check the next validation
	if !objIncentive.IsValid {
		return objIncentive, nil
	}

	// when there is no osl inactive validation then return immediately
	if !validator.Incentive.OslInactiveValidation {
		return objIncentive, nil
	}

	// validate the osl inactive
	restPl := models.ReqOslStatus{Cif: pl.CIF, ProductCode: pl.Validators.Product}
	oslStatus, err := ref.restReferral.RRGetOslStatus(c, restPl)

	if err != nil || !oslStatus {
		return models.ObjIncentive{}, models.ErrOslInactive
	}

	return objIncentive, nil
}

func (rcUc *referralUseCase) UReferralCIFValidate(c echo.Context, cif string) (models.ReferralCodes, error) {

	data, err := rcUc.referralRepo.RGetReferralByCif(c, cif)

	if err != nil {
		return models.ReferralCodes{}, err
	}

	if data.CIF == "" {
		return data, models.ErrCIF
	}

	if data.ReferralCode == "" {
		return data, models.ErrRefCodesNF
	}

	return data, nil
}

func (rcUc *referralUseCase) UGetPrefixActiveCampaignReferral(c echo.Context) (models.PrefixResponse, error) {
	var pv models.PayloadValidator
	var prefixResponse models.PrefixResponse

	pv.TransactionDate = time.Now().Format(models.DateTimeFormat)
	campaignMetadata, err := rcUc.referralRepo.RGetReferralCampaignMetadata(c, pv)

	if err != nil {
		return models.PrefixResponse{}, err
	}

	prefixResponse.Prefix = campaignMetadata.Prefix

	return prefixResponse, nil
}

func (rcUc *referralUseCase) getCampaignIdByPrefix(c echo.Context, prefix string) (int64, error) {
	// get campaign id by prefix
	campaignId, err := rcUc.referralRepo.RGetCampaignId(c, prefix)

	if err != nil {
		return 0, err
	}

	if campaignId == 0 {
		return 0, models.ErrRefCampaginNF
	}

	return campaignId, nil
}
func (rcUc *referralUseCase) UGetHistoryIncentive(c echo.Context, pl models.RequestHistoryIncentive) (models.ResponseHistoryIncentive, error) {
	var historyIncentive models.ResponseHistoryIncentive

	historyIncentive, err := rcUc.referralRepo.RGetHistoryIncentive(c, pl)

	if err != nil {
		return models.ResponseHistoryIncentive{}, models.ErrRefHistoryIncentiveNF
	}

	if historyIncentive.TotalData == 0 || len(*historyIncentive.HistoryIncentiveData) == 0 {
		return models.ResponseHistoryIncentive{}, models.ErrRefHistoryIncentiveNF
	}

	return historyIncentive, nil
}

func (rcUc *referralUseCase) UTotalFriends(c echo.Context, pl models.RequestReferralCodeUser) (models.RespTotalFriends, error) {
	var ttlFriends models.RespTotalFriends

	// Check if a cif already had a referral code
	referral, errRef := rcUc.referralRepo.RGetReferralByCif(c, pl.CIF)

	if errRef != nil {
		return models.RespTotalFriends{}, models.ErrCIF
	}

	// if not exist, return invalid cif
	if referral.ReferralCode == "" {
		return models.RespTotalFriends{}, models.ErrCIF
	}

	ttlFriends, err := rcUc.referralRepo.RTotalFriends(c, pl.CIF)

	if err != nil {
		return models.RespTotalFriends{}, models.ErrCIF
	}

	return ttlFriends, nil
}

func (rcUc *referralUseCase) UFriendsReferral(c echo.Context, pl models.PayloadFriends) ([]models.Friends, error) {

	var referral []models.Friends

	if pl.Limit == 0 {
		pl.Limit = 4
	}

	// get friends refferal by cif
	referral, err := rcUc.referralRepo.RFriendsReferral(c, pl)

	// throw empty data if referral codes not found
	if err != nil {
		return referral, err
	}

	if len(referral) == 0 {
		logger.Make(c).Debug(models.ErrFriends)
		return nil, models.ErrFriends
	}

	return referral, nil
}
