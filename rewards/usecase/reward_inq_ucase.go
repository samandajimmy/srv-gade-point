package usecase

import (
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"

	"time"

	"github.com/labstack/echo"
)

func (rwd *rewardUseCase) Inquiry(c echo.Context, plValidator *models.PayloadValidator) (models.RewardsInquiry, *models.ResponseErrors) {
	var rwdInquiry models.RewardsInquiry
	var rwdResponse []models.RewardResponse
	var respErrors models.ResponseErrors

	// validate trx date
	_, err := time.Parse(models.DateTimeFormatMillisecond, plValidator.TransactionDate)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrTrxDateFormat)
		respErrors.SetTitle(models.ErrTrxDateFormat.Error())

		return rwdInquiry, &respErrors
	}

	// validate the inquiry request, if refId exist
	if plValidator.RefTrx != "" {
		rwdInquiry, err = rwd.rwdTrxRepo.GetByRefID(c, plValidator.RefTrx)

		if (err != nil || rwdInquiry == models.RewardsInquiry{}) {
			logger.Make(c, nil).Debug(models.ErrRefTrxNotFound)
			respErrors.SetTitle(models.ErrRefTrxNotFound.Error())

			return rwdInquiry, &respErrors
		}

		return rwdInquiry, nil
	}

	// get voucher code with promo code
	voucherCode, _, err := rwd.voucherUC.GetVoucherCode(c, plValidator, false)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// check request payload base on cif and promo code
	// get existing reward trx based on cif and phone number
	rwrds, _ := rwd.rwdTrxRepo.GetRewardByPayload(c, *plValidator, voucherCode)

	if len(rwrds) > 0 && voucherCode == nil {
		var rr []models.RewardResponse

		logger.Make(c, nil).Debug(models.ErrMessageRewardTrxAlreadyExists)

		for _, reward := range rwrds {
			// validate each reward
			err = reward.Validators.Validate(plValidator)

			if err != nil {
				logger.Make(c, nil).Debug(err)
				respErrors.AddError(err.Error())

				continue
			}

			// get response reward
			respData, err := rwd.responseReward(c, *reward, plValidator)

			if err != nil {
				logger.Make(c, nil).Debug(err)
				respErrors.AddError(err.Error())

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

		return rwdInquiry, nil
	}

	// check if promoCode is a voucher code
	// if yes then it should call validate voucher
	if voucherCode != nil {
		return rwd.getVoucherReward(c, plValidator, voucherCode)
	}

	// fresh or new reward trx start from here
	// check available campaign
	campaignAvailable, err := rwd.campaignRepo.GetCampaignAvailable(c, *plValidator)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrNoCampaign)
		respErrors.SetTitle(models.ErrNoCampaign.Error())

		return rwdInquiry, &respErrors
	}

	// Logic referral
	// check for referral validate
	isValidate, err := rwd.validateReferralInq(c, plValidator, &respErrors)

	if err != nil {
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	if !isValidate {
		logger.Make(c, nil).Debug(err)

		return rwdInquiry, &respErrors
	}

	// create array reward
	packRewards := rwd.putRewards(c, campaignAvailable)

	for _, reward := range packRewards {
		// validate promo code
		if err = rwd.validatePromoCode(*reward.Tags, reward, plValidator.PromoCode); err != nil {
			logger.Make(c, nil).Debug(err)

			continue
		}

		// validate reward quota
		available, err := rwd.quotaUC.CheckQuota(c, reward, plValidator)

		if !available {
			logger.Make(c, nil).Debug(err)
			respErrors.AddError(err.Error())

			continue
		}

		// validate each reward
		err = reward.Validators.Validate(plValidator)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			respErrors.AddError(err.Error())

			continue
		}

		// get response reward
		rwdResp, err := rwd.responseReward(c, reward, plValidator)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			respErrors.SetTitle(err.Error())

			return rwdInquiry, &respErrors
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
		logger.Make(c, nil).Debug(models.ErrMessageNoRewards)
		respErrors.SetTitle(models.ErrMessageNoRewards.Error())

		return rwdInquiry, &respErrors
	}

	rwdInquiry.Rewards = &rwdResponse

	// if reward greater then one
	if len(*rwdInquiry.Rewards) > 1 {
		rwdInquiry.RefTrx = randRefID(20)
	}

	// insert data to reward transaction
	_, err = rwd.createRewardTrx(c, *plValidator, rwdInquiry)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// if not multi
	if !plValidator.IsMulti {
		rwdInquiry.RefTrx = rwdResponse[0].RefTrx
		rwdResponse[0].RefTrx = ""
	}

	// check referral cant use referrer myself
	if (plValidator.IsMulti) && (plValidator.CIF == plValidator.Referrer) {
		logger.Make(c, nil).Debug(models.ErrSameCifReferrerAndReferral)
		logger.Make(c, nil).Debug(err)
		respErrors.SetTitle(models.ErrSameCifReferrerAndReferral.Error())

		return rwdInquiry, &respErrors
	}

	return rwdInquiry, nil
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
	rwdResp, _ := rwd.responseReward(c, rewardsVoucher[0], plValidator)
	rwdInquiry.RefTrx = rwdResp.RefTrx // nolint
	rwdResp.RewardID = 0               // nolint
	rwdResp.RefTrx = ""                // nolint

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
	requestLogger := logger.Make(c, nil)
	diff := rewardTrx.TimeoutDate.Sub(*rewardTrx.InquiredDate)
	delay := time.Duration(diff.Seconds())

	go func(rwdTrx models.RewardTrx, delay time.Duration) {
		requestLogger.Debug("Store job to background for ref ID: " + rwdTrx.RefID)
		time.Sleep(delay * time.Second)
		requestLogger.Debug("Start to make ref ID: " + rwdTrx.RefID + " expired!")
		rwd.rwdTrxRepo.RewardTrxTimeout(rwdTrx)
	}(rewardTrx, delay)
}
