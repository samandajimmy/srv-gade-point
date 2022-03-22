package usecase

import (
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"math/rand"
	"time"

	"github.com/labstack/echo"
)

const (
	letters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lengthCode = 10
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

func (rcUc *referralUseCase) CreateReferralCodes(c echo.Context, requestReferralCodes models.RequestCreateReferral) (models.ReferralCodes, error) {

	var payload models.ReferralCodes
	var pv models.PayloadValidator

	// Check if a cif already had a referral code
	referral, err := rcUc.referralRepo.GetReferralByCif(c, requestReferralCodes.CIF)

	if err != nil {
		return models.ReferralCodes{}, err
	}

	pv.PromoCode = referral.ReferralCode
	pv.TransactionDate = time.Now().Format(models.DateTimeFormat)

	// check if referralCode has reach max trx day/month
	cr := rcUc.campaignRepo.GetReferralCampaign(c, pv)
	_, err = rcUc.validateReferralIncentive(c, cr, pv)

	if err != nil {
		return models.ReferralCodes{}, err
	}

	// if exist, return existing referral code
	if referral.ReferralCode != "" {
		logger.Make(c, nil).Info(models.DynamicErr(models.CifRefCodeExisted, requestReferralCodes.CIF))

		return referral, models.DynamicErr(models.CifRefCodeExisted, requestReferralCodes.CIF)
	}

	// if not, generate referral code, getting campaign id first
	campaignId, err := rcUc.getCampaignIdByPrefix(c, requestReferralCodes.Prefix)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	// map campaign id, generate referral code according to requirements (with prefix)
	payload.CampaignId = campaignId
	payload.ReferralCode = rcUc.generateRefCode(requestReferralCodes.Prefix)
	payload.CIF = requestReferralCodes.CIF

	// create referral code data
	result, err := rcUc.referralRepo.CreateReferral(c, payload)

	if err != nil {
		return models.ReferralCodes{}, models.ErrGenerateRefCodeFailed
	}

	return result, nil
}

func (rcUc *referralUseCase) GetReferralCodes(c echo.Context, requestGetReferral models.RequestReferralCodeUser) (models.ReferralCodeUser, error) {

	var result models.ReferralCodeUser

	// get referral codes data by cif
	referral, err := rcUc.referralRepo.GetReferralByCif(c, requestGetReferral.CIF)

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
	campaignId, err := rcUc.referralRepo.GetCampaignId(c, prefix)

	if err != nil {
		return 0, err
	}

	if campaignId == 0 {
		logger.Make(c, nil).Debug(models.ErrRefCampaginNF)

		return 0, models.ErrRefCampaginNF
	}

	return campaignId, nil
}

func (rcUc *referralUseCase) generateRefCode(prefix string) string {

	lettersUsed := []rune(letters)
	prefixLength := len(prefix)

	var stringGeneratedLength = lengthCode - prefixLength

	s := make([]rune, stringGeneratedLength)
	for i := range s {
		s[i] = lettersUsed[rand.Intn(len(lettersUsed))]
	}

	return prefix + string(s)
}

func (rcUc *referralUseCase) UReferralCIFValidate(c echo.Context, cif string) (models.ReferralCodes, error) {

	data, err := rcUc.referralRepo.GetReferralByCif(c, cif)

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

func (rcUc *referralUseCase) validateReferralIncentive(c echo.Context, cp *[]*models.Campaign, pv models.PayloadValidator) (models.SumIncentive, error) {

	if cp == nil {
		return models.SumIncentive{}, nil
	}

	camps := *cp
	campaign := *camps[0]

	sumIncentive, err := rcUc.UValidateReferrer(c, pv, &campaign)

	if err != nil {
		return models.SumIncentive{}, err
	}

	if !sumIncentive.ValidPerDay {
		logger.Make(c, nil).Error(models.DynamicErr(models.DynErrIncDayExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerDay)))

		return models.SumIncentive{}, models.DynamicErr(models.DynErrIncDayExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerDay))
	}

	if !sumIncentive.ValidPerMonth {
		logger.Make(c, nil).Error(models.DynamicErr(models.DynErrIncMonthExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerMonth)))

		return models.SumIncentive{}, models.DynamicErr(models.DynErrIncMonthExceeded, helper.FloatToString(sumIncentive.Reward.Validators.Incentive.MaxPerMonth))
	}

	return sumIncentive, nil
}

func (rcUc *referralUseCase) UGetHistoryIncentive(c echo.Context, pl models.RequestHistoryIncentive) ([]models.ResponseHistoryIncentive, error) {
	var historyIncentive []models.ResponseHistoryIncentive

	historyIncentive, err := rcUc.referralRepo.RGetHistoryIncentive(c, pl.RefCif)

	if err != nil {
		logger.Make(c, nil).Error(models.ErrRefHistoryIncentiveNF)

		return []models.ResponseHistoryIncentive{}, models.ErrRefHistoryIncentiveNF
	}

	length := len(historyIncentive)

	if length == 0 {
		logger.Make(c, nil).Error(models.ErrRefHistoryIncentiveNF)

		return []models.ResponseHistoryIncentive{}, models.ErrRefHistoryIncentiveNF
	}

	return historyIncentive, nil
}
