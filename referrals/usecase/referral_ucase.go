package usecase

import (
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"github.com/labstack/echo"
)

type referralUseCase struct {
	referralRepo referrals.RefRepository
	campaignRepo campaigns.CRepository
}

// NewReferralUseCase will create new an rewardtrxUseCase object representation of referrals.RefUseCase interface
func NewReferralUseCase(
	refRepo referrals.RefRepository,
	campaignRepo campaigns.CRepository,
) referrals.RefUseCase {
	return &referralUseCase{
		referralRepo: refRepo,
		campaignRepo: campaignRepo,
	}
}

func (ref *referralUseCase) UPostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) ([]models.CoreTrxResponse, error) {
	var responseData []models.CoreTrxResponse
	err := ref.referralRepo.RPostCoreTrx(c, coreTrx)

	if len(coreTrx) == 0 {
		return nil, models.ErrCoreTrxEmpty
	}

	if err != nil {
		return nil, models.ErrCoreTrxFailed
	}

	return responseData, nil
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
		logger.Make(c, nil).Debug(err)

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

func (rcUc *referralUseCase) UGetReferralCodes(c echo.Context, requestGetReferral models.RequestReferralCodeUser) (models.ReferralCodeUser, error) {

	var result models.ReferralCodeUser

	// get referral codes data by cif
	referral, err := rcUc.referralRepo.RGetReferralByCif(c, requestGetReferral.CIF)

	// throw empty data if referral codes not found
	if err != nil {
		return models.ReferralCodeUser{}, err
	}

	if referral.ReferralCode == "" {
		logger.Make(c, nil).Debug(models.ErrRefCodesNF)

		return models.ReferralCodeUser{}, models.ErrRefCodesNF
	}

	// if data referral codes found send result
	result.ReferralCode = referral.ReferralCode

	return result, nil
}

func (ref *referralUseCase) UValidateReferrer(c echo.Context, pl models.PayloadValidator, campaign *models.Campaign) (models.SumIncentive, error) {
	var sumIncentive models.SumIncentive
	var err error

	// get ref trx based on the active ref campaign
	// and reward data with used ref code
	if campaign == nil {
		camps := ref.campaignRepo.GetReferralCampaign(c, pl)
		c := *camps
		campaign = c[0]
	}

	for _, reward := range *campaign.Rewards {
		if reward.Validators.Target != models.RefTargetReferrer {
			continue
		}

		// sum the incentive based on the ref trx (per day and per month)
		sumIncentive, err = ref.referralRepo.RSumRefIncentive(c, pl.PromoCode, reward)

		if err != nil {
			continue
		}

		if (sumIncentive.Reward != models.Reward{}) {
			break
		}

	}

	// validate the max per day and per month
	validator := sumIncentive.Reward.Validators
	validator.Incentive.ValidateMaxIncentive(&sumIncentive)

	return sumIncentive, nil
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

// func (rcUc *referralUseCase) validateReferralIncentive(c echo.Context, cp *[]*models.Campaign, pv models.PayloadValidator) (models.SumIncentive, error) {

// 	if cp == nil {
// 		return models.SumIncentive{}, nil
// 	}

// 	camps := *cp
// 	campaign := *camps[0]

// 	sumIncentive, err := rcUc.UValidateReferrer(c, pv, &campaign)

// 	if err != nil {
// 		return models.SumIncentive{}, err
// 	}

// 	if !sumIncentive.ValidPerDay {
// 		logger.Make(c, nil).Error(models.DynamicErr(models.DynErrIncDayExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerDay)))

// 		return models.SumIncentive{}, models.DynamicErr(models.DynErrIncDayExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerDay))
// 	}

// 	if !sumIncentive.ValidPerMonth {
// 		logger.Make(c, nil).Error(models.DynamicErr(models.DynErrIncMonthExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerMonth)))

// 		return models.SumIncentive{}, models.DynamicErr(models.DynErrIncMonthExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerMonth))
// 	}

// 	return sumIncentive, nil
// }
