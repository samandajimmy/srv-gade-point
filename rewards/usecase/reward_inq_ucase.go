package usecase

import (
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"strconv"

	"time"

	"github.com/labstack/echo"
)

func (rwd *rewardUseCase) Inquiry(c echo.Context, plValidator *models.PayloadValidator) (models.RewardsInquiry, *models.ResponseErrors) {
	var rwdInquiry models.RewardsInquiry
	var rwdResponse []models.RewardResponse
	var respErrors models.ResponseErrors

	cp := models.CampaignPack{}
	// get code type promo, voucher or referral
	plValidator.CodeType = rwd.getCodeType(c, plValidator, &cp)
	// get existing reward transaction if exist
	existingTrx := rwd.getExistingTrx(c, plValidator)

	if existingTrx != nil {
		return *existingTrx, &respErrors
	}

	// check if promoCode is a voucher code
	// if yes then it should call get voucher reward
	if plValidator.CodeType == models.CodeTypeVoucher {
		return rwd.getVoucherReward(c, plValidator, cp.VoucherPack)
	}

	// check if promoCode type is a referral code. if yes,
	// we need to validate the referral promo campaign
	err := rwd.validatePromoReferral(c, plValidator, cp.ReferralPack)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// fresh or new reward trx start from here
	// get rewards availables
	idxStr := ""
	packRewards, err := rwd.getAvailableRewards(c, *plValidator, &cp)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	for idx, reward := range packRewards {
		idxStr = strconv.Itoa(idx)
		// validate promo code
		if err = rwd.validatePromoCode(*reward.Tags, reward, plValidator.PromoCode); err != nil {
			continue
		}

		// validate reward quota
		available, err := rwd.quotaUC.CheckQuota(c, reward, plValidator)

		if !available {
			respErrors.AddError(
				models.DynamicErr(models.ErrDynRewardValidation, idxStr, err.Error()).Error())

			continue
		}

		// get response reward
		rwdResp, err := rwd.responseReward(c, reward, plValidator)

		if err != nil {
			respErrors.AddError(
				models.DynamicErr(models.ErrDynRewardValidation, idxStr, err.Error()).Error())

			continue
		}

		if rwdResp != nil {
			rwdResponse = append(rwdResponse, *rwdResp)
		}

		// update reward quota
		_ = rwd.quotaUC.UpdateReduceQuota(c, reward.ID)

		// if not multi
		if !plValidator.IsMulti {
			break
		}
	}

	// if no reward found
	if len(rwdResponse) == 0 {
		respErrors.SetTitle(models.ErrMessageNoRewards.Error())

		return rwdInquiry, &respErrors
	}

	rwdInquiry.Rewards = &rwdResponse
	respErrors = models.ResponseErrors{}
	rwdInquiry.RefTrx = rwdResponse[0].RefTrx

	// if reward greater then one
	if len(*rwdInquiry.Rewards) > 1 {
		rwdInquiry.RefTrx = rwd.rewardRepo.RGetRandomId(20)
	}

	// insert data to reward transaction
	_, err = rwd.createRewardTrx(c, *plValidator, rwdInquiry)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// if not multi
	if !plValidator.IsMulti {
		rwdInquiry.RefTrx = rwdResponse[0].RefTrx
		rwdResponse[0].RefTrx = ""
	}

	return rwdInquiry, &respErrors
}

func (rwd *rewardUseCase) validatePromoReferral(c echo.Context, pl *models.PayloadValidator,
	campaigns *[]*models.Campaign) error {

	if pl.CodeType != models.CodeTypeReferral {
		return nil
	}

	// validate referral codes
	camps := *campaigns
	campaignReferral := models.CampaignReferral{
		Campaign:    *camps[0],
		CifReferrer: pl.Referrer,
	}
	// get data incentive of used referral code
	incentive, err := rwd.referralUs.UValidateReferrer(c, *pl, &campaignReferral)

	if err != nil {
		return err
	}

	// check incentive data valid (perDay and perMonth) or not
	if !incentive.ValidPerDay || !incentive.ValidPerMonth {
		return models.ErrRefTrxIncentive
	}

	// check incentive data valid or not
	if !incentive.IsValid {
		return models.ErrRefTrxExceeded
	}

	// validate referral code only one CIF per product code
	isValidRefCode := rwd.referralTrxRepo.RGetCountRefTrxByCIF(c, pl.CIF, pl.Validators.Product, pl.PromoCode)

	// check incentive data valid or not
	if !isValidRefCode {
		return models.ErrRefTrxExceededPC
	}

	return nil
}

func (rwd *rewardUseCase) getVoucherReward(c echo.Context, plValidator *models.PayloadValidator,
	vc *models.VoucherCode) (models.RewardsInquiry, *models.ResponseErrors) {

	var rwdInquiry models.RewardsInquiry
	var rwdResponse []models.RewardResponse
	var respErrors models.ResponseErrors

	rewardsVoucher, err := rwd.voucherUC.VoucherValidate(c, plValidator, vc)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// get response reward
	rwdResp, _ := rwd.responseInquiry(c, rewardsVoucher[0], plValidator)
	rwdInquiry.RefTrx = rwdResp.RefTrx // nolint
	rwdResp.RewardID = 0               // nolint

	// nolint
	if rwdResp != nil {
		rwdResponse = append(rwdResponse, *rwdResp)
	}

	rwdInquiry.Rewards = &rwdResponse

	// insert data to reward transaction
	_, err = rwd.createRewardTrx(c, *plValidator, rwdInquiry)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	return rwdInquiry, &respErrors
}

func (rwd *rewardUseCase) timeoutTrxJob(c echo.Context, rewardTrx models.RewardTrx) {
	requestLogger := logger.Make(c)
	diff := rewardTrx.TimeoutDate.Sub(*rewardTrx.InquiredDate)
	delay := time.Duration(diff.Seconds())

	go func(rwdTrx models.RewardTrx, delay time.Duration) {
		requestLogger.Debug("Store job to background for ref ID: " + rwdTrx.RefID)
		time.Sleep(delay * time.Second)
		requestLogger.Debug("Start to make ref ID: " + rwdTrx.RefID + " expired!")
		rwd.rwdTrxRepo.RewardTrxTimeout(rwdTrx)
	}(rewardTrx, delay)
}

func (rwd *rewardUseCase) getAvailableRewards(c echo.Context, pl models.PayloadValidator,
	cp *models.CampaignPack) ([]models.Reward, error) {
	var rewards []models.Reward

	if cp.ReferralPack != nil {
		return rwd.putRewards(c, *cp.ReferralPack), nil
	}

	// check available campaign
	campaignAvailable, err := rwd.campaignRepo.GetCampaignAvailable(c, pl)

	if err != nil {
		return rewards, models.ErrNoCampaign
	}

	cp.ReferralPack = &campaignAvailable
	// create array reward
	rewards = rwd.putRewards(c, campaignAvailable)

	return rewards, nil
}

func (rwd *rewardUseCase) responseReward(c echo.Context, reward models.Reward,
	plValidator *models.PayloadValidator) (*models.RewardResponse, error) {
	// validate each reward
	err := reward.Validators.Validate(plValidator)

	if err != nil {
		return nil, err
	}

	rwdResp, err := rwd.responseInquiry(c, reward, plValidator)

	if err != nil {
		return nil, err
	}

	return rwdResp, nil
}

func (rwd *rewardUseCase) getCodeType(c echo.Context, pl *models.PayloadValidator,
	cp *models.CampaignPack) string {

	// check if the promo code is a voucher or not
	// get voucher code with promo code
	vc, _, _ := rwd.voucherUC.GetVoucherCode(c, pl, false)

	if vc != nil {
		cp.VoucherPack = vc

		return models.CodeTypeVoucher
	}

	// check if the promo code is a referral or not
	// get campaign referral by referral code
	cr := rwd.campaignRepo.GetReferralCampaign(c, *pl)

	if cr != nil {
		var campaign []*models.Campaign
		campaign = append(campaign, &cr.Campaign)
		cp.ReferralPack = &campaign
		pl.Referrer = cr.CifReferrer

		return models.CodeTypeReferral
	}

	return models.CodeTypePromo
}

func (rwd *rewardUseCase) getExistingTrx(c echo.Context, plValidator *models.PayloadValidator) *models.RewardsInquiry {
	var rwdInquiry *models.RewardsInquiry

	// validate the inquiry request, if refId exist
	// then get directly from existing reward trx by refId
	rwdInquiry = rwd.getInquiredTrxByRefTrx(c, plValidator)

	if rwdInquiry != nil {
		return rwdInquiry
	}

	// check request payload base on cif and promo code
	// get existing reward trx based on cif and phone number
	rwdInquiry = rwd.getInquiredTrxByPayload(c, plValidator)

	if rwdInquiry != nil {
		return rwdInquiry
	}

	return nil
}

func (rwd *rewardUseCase) getInquiredTrxByRefTrx(c echo.Context, plValidator *models.PayloadValidator) *models.RewardsInquiry {
	var rwdInquiry models.RewardsInquiry

	if plValidator.RefTrx == "" {
		return nil
	}

	rwdInquiry, err := rwd.rwdTrxRepo.GetByRefID(c, plValidator.RefTrx)

	if (err != nil || rwdInquiry == models.RewardsInquiry{}) {
		logger.Make(c).Debug(models.ErrRefTrxNotFound)

		return nil
	}

	return &rwdInquiry
}

func (rwd *rewardUseCase) getInquiredTrxByPayload(c echo.Context, plValidator *models.PayloadValidator) *models.RewardsInquiry {
	var rwdInquiry models.RewardsInquiry
	var rr []models.RewardResponse

	if plValidator.CodeType == models.CodeTypeVoucher {
		return nil
	}

	rwrds, err := rwd.rwdTrxRepo.GetRewardByPayload(c, *plValidator, nil)

	if err != nil || len(rwrds) == 0 {
		return nil
	}

	logger.Make(c).Debug(models.ErrMessageRewardTrxAlreadyExists)

	for _, reward := range rwrds {
		// get response reward
		respData, err := rwd.responseReward(c, *reward, plValidator)

		if err != nil {
			continue
		}

		rr = append(rr, *respData)
	}

	rwdInquiry.Rewards = &rr

	if rwrds[0].RootRefID != "" {
		rwdInquiry.RefTrx = rwrds[0].RootRefID
	}

	// if not multi
	if !plValidator.IsMulti {
		rwdInquiry.RefTrx = rr[0].RefTrx
		rr[0].RefTrx = ""
	}

	return &rwdInquiry
}
