package usecase

import (
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/rewardtrxs"
	"gade/srv-gade-point/tags"
	"gade/srv-gade-point/vouchercodes"
	"gade/srv-gade-point/vouchers"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type rewardUseCase struct {
	rewardRepo      rewards.Repository
	campaignRepo    campaigns.Repository
	tagUC           tags.UseCase
	quotaUC         quotas.UseCase
	voucherUC       vouchers.UseCase
	voucherCodeRepo vouchercodes.Repository
	rwdTrxRepo      rewardtrxs.Repository
}

// NewRewardUseCase will create new an rewardUseCase object representation of rewards.UseCase interface
func NewRewardUseCase(
	rwdRepo rewards.Repository,
	campaignRepo campaigns.Repository,
	tagUC tags.UseCase,
	quotaUC quotas.UseCase,
	voucherUC vouchers.UseCase,
	voucherCodeRepo vouchercodes.Repository,
	rwdTrxRepo rewardtrxs.Repository,
) rewards.UseCase {
	return &rewardUseCase{
		rewardRepo:      rwdRepo,
		campaignRepo:    campaignRepo,
		tagUC:           tagUC,
		quotaUC:         quotaUC,
		voucherUC:       voucherUC,
		voucherCodeRepo: voucherCodeRepo,
		rwdTrxRepo:      rwdTrxRepo,
	}
}

func (rwd *rewardUseCase) CreateReward(c echo.Context, reward *models.Reward, campaign *models.Campaign) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	reward.Campaign = campaign
	err := rwd.rewardRepo.CreateReward(c, reward, campaign.ID)

	if err != nil {
		requestLogger.Debug(models.ErrRewardFailed)

		return models.ErrRewardFailed
	}

	// create array quotas
	if reward.Quotas != nil {
		for _, quota := range *reward.Quotas {
			err = rwd.quotaUC.Create(c, &quota, reward)

			if err != nil {
				break
			}
		}

		if err != nil {
			_ = rwd.quotaUC.DeleteByReward(c, reward.ID)
			requestLogger.Debug(models.ErrCreateQuotasFailed)

			return models.ErrCreateQuotasFailed
		}
	}

	// create array tags
	if reward.Tags != nil {
		for _, tag := range *reward.Tags {
			err = rwd.tagUC.CreateTag(c, &tag, reward.ID)

			if err != nil {
				break
			}

			err = rwd.rewardRepo.CreateRewardTag(c, &tag, reward.ID)

			if err != nil {
				break
			}
		}

		if err != nil {
			_ = rwd.rewardRepo.DeleteRewardTag(c, reward.ID)
			requestLogger.Debug(models.ErrCreateTagsFailed)

			return models.ErrCreateTagsFailed
		}
	}

	return nil
}

func (rwd *rewardUseCase) DeleteByCampaign(c echo.Context, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwd.rewardRepo.DeleteByCampaign(c, campaignID)

	if err != nil {
		requestLogger.Debug(models.ErrDelRewardFailed)

		return models.ErrDelRewardFailed
	}

	return nil
}

func (rwd *rewardUseCase) Inquiry(c echo.Context, plValidator *models.PayloadValidator) (models.RewardsInquiry, *models.ResponseErrors) {
	var rwdInquiry models.RewardsInquiry
	var rwdResponse []models.RewardResponse
	var voucherCode *models.VoucherCode
	var respErrors models.ResponseErrors
	var rewardSelected models.Reward
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	trxDate, err := time.Parse(models.DateTimeFormatMillisecond, plValidator.TransactionDate)

	if err != nil {
		requestLogger.Debug(models.ErrTrxDateFormat)
		respErrors.SetTitle(models.ErrTrxDateFormat.Error())

		return rwdInquiry, &respErrors
	}

	// validate the inquiry request, if refId exist
	if plValidator.RefTrx != "" {
		rwdInquiry, err = rwd.rwdTrxRepo.GetByRefID(c, plValidator.RefTrx)

		if err != nil {
			requestLogger.Debug(models.ErrRefTrxNotFound)
			respErrors.SetTitle(models.ErrRefTrxNotFound.Error())

			return rwdInquiry, &respErrors
		}

		return rwdInquiry, nil
	}

	// check request payload base on transactionDate
	rwdTrxInquiry, err := rwd.rwdTrxRepo.CheckByTransactionDate(c, *plValidator)

	if rwdTrxInquiry != nil {
		requestLogger.Debug(models.ErrMessageRewardTrxAlreadyExists)

		return *rwdTrxInquiry, nil
	}

	// check available campaign
	campaigns, err := rwd.campaignRepo.GetCampaignAvailable(c, trxDate.Format(models.DateTimeFormat))

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)
		respErrors.SetTitle(models.ErrNoCampaign.Error())

		return rwdInquiry, &respErrors
	}

	// create array rewards
	rewards := rwd.putRewards(c, campaigns)

	for _, reward := range rewards {
		var rwdResp models.RewardResponse
		rewardLogger := logger.GetRequestLogger(c, reward.Validators)

		// validate promo code
		if err = rwd.validatePromoCode(*reward.Tags, reward.PromoCode, plValidator.PromoCode); err != nil {
			rewardLogger.Debug(err)

			continue
		}

		// validate reward quota
		available, err := rwd.quotaUC.CheckQuota(c, reward, plValidator)

		if available == false {
			rewardLogger.Debug(err)
			respErrors.AddError(err.Error())

			continue
		}

		// validate each reward
		if err = reward.Validators.Validate(plValidator); err != nil {
			rewardLogger.Debug(err)
			respErrors.AddError(err.Error())

			continue
		}

		// get the rewards value/benefit
		rwdValue, _ := reward.Validators.GetRewardValue(plValidator)

		// check rewards voucher if any
		rwdVoucher, _ := reward.Validators.GetVoucherResult()

		if rwdVoucher != 0 {
			plVoucherBuy := &models.PayloadVoucherBuy{
				CIF:       plValidator.CIF,
				VoucherID: strconv.FormatInt(rwdVoucher, 10),
			}

			voucherCode, err = rwd.voucherUC.VoucherGive(c, plVoucherBuy)

			if err != nil {
				requestLogger.Debug(models.ErrVoucherUnavailable)
				respErrors.AddError(models.ErrVoucherUnavailable.Error())

				continue
			}

			rwdResp.VoucherName = voucherCode.Voucher.Name
			rwdValue = 0 // if voucher reward is exist then reward value should be nil
		}

		// populate reward response
		rwdResp.Type = reward.GetRewardTypeText()
		rwdResp.Value = rwdValue
		rwdResp.JournalAccount = reward.JournalAccount
		rwdResponse = append(rwdResponse, rwdResp)
		rewardSelected = reward
	}

	if len(rwdResponse) == 0 {
		requestLogger.Debug(models.ErrMessageNoRewards)
		respErrors.SetTitle(models.ErrMessageNoRewards.Error())

		return rwdInquiry, &respErrors
	}

	// insert data to reward history
	rewardTrx, err := rwd.createRewardTrx(c, *plValidator, rewardSelected.ID, rwdResponse)

	if err != nil {
		requestLogger.Debug(err)
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// update voucher code ref_id
	err = rwd.voucherCodeRepo.UpdateVoucherCodeRefID(c, voucherCode, rewardTrx.RefID)

	if err != nil {
		requestLogger.Debug(err)
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// generate an unique ref ID
	rwdInquiry.RefTrx = rewardTrx.RefID
	rwdInquiry.Rewards = &rwdResponse

	// update reward quota
	rwd.quotaUC.UpdateReduceQuota(c, rewardSelected.ID)

	return rwdInquiry, nil
}

func (rwd *rewardUseCase) Payment(c echo.Context, rwdPayment *models.RewardPayment) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check available reward transaction based in CIF and ref_id
	rewardtrx, err := rwd.rwdTrxRepo.CheckTrx(c, rwdPayment.CIF, rwdPayment.RefTrx)

	if err != nil {
		requestLogger.Debug(models.ErrRefTrxNotFound)

		return models.ErrRefTrxNotFound
	}

	if rwdPayment.RefCore == "" {
		// rejected
		// update voucher code
		rwd.voucherCodeRepo.UpdateVoucherCodeRejected(c, rwdPayment.RefTrx)

		// update reward trx
		rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, models.RewardTrxRejected)

		// update add reward quota
		rwd.quotaUC.UpdateAddQuota(c, *rewardtrx.RewardID)

		// update reward quota
	} else {
		// succeeded
		// update reward trx
		rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, models.RewardTrxSucceeded)

	}

	return nil
}

func (rwd *rewardUseCase) CheckTransaction(c echo.Context, rwdPayment *models.RewardPayment) (models.RewardTrxResponse, error) {
	var responseData models.RewardTrxResponse
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check available reward transaction based in ref_id
	rewardtrx, err := rwd.rwdTrxRepo.CheckRefID(c, rwdPayment.RefTrx)

	if err != nil {
		requestLogger.Debug(models.ErrRefTrxNotFound)

		return responseData, models.ErrRefTrxNotFound
	}

	responseData.StatusCode = rewardtrx.Status
	responseData.Status = rewardtrx.GetstatusRewardTrxText()

	return responseData, nil
}

func (rwd *rewardUseCase) RefreshTrx() {
	now := time.Now()
	// update trx that should be timeout
	err := rwd.rwdTrxRepo.UpdateTimeoutTrx()

	if err != nil {
		logrus.Debug(err)
	}

	// get trx that need to be timeout later
	rewardTrx, err := rwd.rwdTrxRepo.GetInquiredTrx()

	if err != nil {
		logrus.Debug(err)
	}

	for _, rwdTrx := range rewardTrx {
		diff := rwdTrx.TimeoutDate.Sub(now)
		delay := time.Duration(diff.Seconds())

		go func(rwdTrx models.RewardTrx, delay time.Duration) {
			logrus.Debug("Store job to background for ref ID: " + rwdTrx.RefID)
			time.Sleep(delay * time.Second)
			logrus.Debug("Start to make ref ID: " + rwdTrx.RefID + " expired!")
			rwd.rwdTrxRepo.RewardTrxTimeout(rwdTrx)
		}(rwdTrx, delay)

	}
}

func (rwd *rewardUseCase) putRewards(c echo.Context, campaigns []*models.Campaign) []models.Reward {
	var rewards []models.Reward

	for _, campaign := range campaigns {
		for _, reward := range *campaign.Rewards {
			rwd.rewardRepo.GetRewardTags(c, &reward)
			reward.Campaign = campaign
			rewards = append(rewards, reward)
		}
	}

	return rewards
}

func (rwd *rewardUseCase) validatePromoCode(tags []models.Tag, validPC, promoCode string) error {
	if promoCode == validPC {
		return nil
	}

	for _, tag := range tags {
		if promoCode == tag.Name {
			return nil
		}
	}

	return models.ErrPromoCode
}

func (rwd *rewardUseCase) createRewardTrx(c echo.Context, plValidator models.PayloadValidator, rewardID int64, rwdResponse []models.RewardResponse) (*models.RewardTrx, error) {
	rewardTrx, err := rwd.rwdTrxRepo.Create(c, plValidator, rewardID, rwdResponse)

	if err != nil {
		return nil, err
	}

	rwd.timeoutTrxJob(c, rewardTrx)

	return &rewardTrx, nil
}

func (rwd *rewardUseCase) timeoutTrxJob(c echo.Context, rewardTrx models.RewardTrx) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	diff := rewardTrx.TimeoutDate.Sub(*rewardTrx.InquiredDate)
	delay := time.Duration(diff.Seconds())

	go func(rwdTrx models.RewardTrx, delay time.Duration) {
		requestLogger.Debug("Store job to background for ref ID: " + rwdTrx.RefID)
		time.Sleep(delay * time.Second)
		requestLogger.Debug("Start to make ref ID: " + rwdTrx.RefID + " expired!")
		rwd.rwdTrxRepo.RewardTrxTimeout(rwdTrx)
	}(rewardTrx, delay)
}
