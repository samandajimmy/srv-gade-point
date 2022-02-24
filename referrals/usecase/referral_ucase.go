package usecase

import (
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"math/rand"

	"github.com/labstack/echo"
)

const (
	letters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lengthCode = 10
)

type referralUseCase struct {
	referralRepo referrals.Repository
}

// NewReferralUseCase will create new an rewardtrxUseCase object representation of referrals.UseCase interface
func NewReferralUseCase(refRepo referrals.Repository) referrals.UseCase {
	return &referralUseCase{
		referralRepo: refRepo,
	}
}

func (rcUc *referralUseCase) CreateReferralCodes(c echo.Context, requestReferralCodes models.RequestCreateReferral) (models.ReferralCodes, error) {

	var payload models.ReferralCodes

	// Check if a cif already had a referral code
	referral, err := rcUc.referralRepo.GetReferralByCif(c, requestReferralCodes.CIF)

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
