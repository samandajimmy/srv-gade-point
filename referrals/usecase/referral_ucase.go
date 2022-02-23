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

func (rcUc *referralUseCase) CreateReferralCodes(c echo.Context, referralCodes models.ReferralCodes) (models.ReferralCodes, error) {

	var payload models.ReferralCodes
	payload.CIF = referralCodes.CIF
	payload.Prefix = referralCodes.Prefix

	// Check if a cif already had a referral code
	referral, _ := rcUc.referralRepo.GetReferralCodesByCif(c, payload)

	// if exist, return existing referral code
	if referral.ReferralCode != "" {
		logger.Make(c, nil).Info(models.DynamicErr(models.AlreadyHadRefCodes, referralCodes.CIF))

		return referral, nil
	}

	// if not, generate referral code, getting campaign id first
	campaignId, err := rcUc.GetCampaignByPrefix(c, referralCodes.Prefix)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	// map campaign id, generate referral code according to requirements (with prefix)
	payload.CampaignId = campaignId
	payload.ReferralCode = GenerateReferalCode(payload.Prefix)

	// create referral code data
	result, err := rcUc.referralRepo.CreateReferralCodes(c, payload)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrReferralCodeFailed)

		return models.ReferralCodes{}, models.ErrReferralCodeFailed
	}

	return result, nil
}

func (rcUc *referralUseCase) GetCampaignByPrefix(c echo.Context, prefix string) (int64, error) {

	// get campaign id by prefix
	campaignId, err := rcUc.referralRepo.GetCampaignByPrefix(c, prefix)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	return campaignId, nil
}

func GenerateReferalCode(prefix string) string {

	lettersUsed := []rune(letters)
	prefixLength := len(prefix)

	var stringGeneratedLength = lengthCode - prefixLength

	s := make([]rune, stringGeneratedLength)
	for i := range s {
		s[i] = lettersUsed[rand.Intn(len(lettersUsed))]
	}

	return prefix + string(s)
}
