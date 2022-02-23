package usecase

import (
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

func (rcUc *referralUseCase) CreateReferralCodes(c echo.Context, referralCodes *models.ReferralCodes) error {

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// Check if a cif already had a referral code
	_ = rcUc.referralRepo.GetReferralCodesByCif(c, referralCodes)

	// if exist, return existing referral code
	if referralCodes.ReferralCode != "" {
		requestLogger.Info(models.DynamicErr(models.AlreadyHadRefCodes, referralCodes.CIF))

		return nil
	}

	// if not, generate referral code, getting campaign id first
	campaignId, err := rcUc.GetCampaignByPrefix(c, referralCodes.Prefix)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	// map campaign id, generate referral code according to requirements (with prefix)
	referralCodes.CampaignId = campaignId
	referralCodes.ReferralCode = GenerateReferalCode(referralCodes.Prefix)

	// create referral code data
	err = rcUc.referralRepo.CreateReferralCodes(c, referralCodes)

	if err != nil {
		requestLogger.Debug(models.ErrReferralCodeFailed)

		return models.ErrReferralCodeFailed
	}

	return nil
}

func (rcUc *referralUseCase) GetCampaignByPrefix(c echo.Context, prefix string) (int64, error) {

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// get campaign id by prefix
	campaignId, err := rcUc.referralRepo.GetCampaignByPrefix(c, prefix)

	if err != nil {
		requestLogger.Debug(err)

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
