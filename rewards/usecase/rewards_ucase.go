package usecase

import (
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/referraltrxs"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/rewardtrxs"
	"gade/srv-gade-point/tags"
	"gade/srv-gade-point/vouchercodes"
	"gade/srv-gade-point/vouchers"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	referralTrxRepo referraltrxs.Repository
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
	referralTrxRepo referraltrxs.Repository,
) rewards.UseCase {
	return &rewardUseCase{
		rewardRepo:      rwdRepo,
		campaignRepo:    campaignRepo,
		tagUC:           tagUC,
		quotaUC:         quotaUC,
		voucherUC:       voucherUC,
		voucherCodeRepo: voucherCodeRepo,
		rwdTrxRepo:      rwdTrxRepo,
		referralTrxRepo: referralTrxRepo,
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
	var respErrors models.ResponseErrors
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// validate trx date
	_, err := time.Parse(models.DateTimeFormatMillisecond, plValidator.TransactionDate)

	if err != nil {
		requestLogger.Debug(models.ErrTrxDateFormat)
		respErrors.SetTitle(models.ErrTrxDateFormat.Error())

		return rwdInquiry, &respErrors
	}

	// validate the inquiry request, if refId exist
	if plValidator.RefTrx != "" {
		rwdInquiry, err = rwd.rwdTrxRepo.GetByRefID(c, plValidator.RefTrx)

		if (err != nil || rwdInquiry == models.RewardsInquiry{}) {
			requestLogger.Debug(models.ErrRefTrxNotFound)
			respErrors.SetTitle(models.ErrRefTrxNotFound.Error())

			return rwdInquiry, &respErrors
		}

		return rwdInquiry, nil
	}

	// check request payload base on cif and promo code
	// get existing reward trx based on cif and phone number
	rwrds, err := rwd.rwdTrxRepo.GetRewardByPayload(c, *plValidator)

	if len(rwrds) > 0 {
		var rr []models.RewardResponse

		requestLogger.Debug(models.ErrMessageRewardTrxAlreadyExists)

		for _, reward := range rwrds {
			// validate each reward
			err = reward.Validators.Validate(plValidator)

			if err != nil {
				requestLogger.Debug(err)
				respErrors.AddError(err.Error())

				continue
			}

			// get response reward
			respData, err := rwd.responseReward(c, *reward, plValidator)

			if err != nil {
				requestLogger.Debug(err)
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
		if plValidator.IsMulti == false {
			rwdInquiry.RefTrx = rr[0].RefTrx
			rr[0].RefTrx = ""
		}

		return rwdInquiry, nil
	}

	// check if promoCode is a voucher code
	// if yes then it should call validate voucher
	voucherCode, _, _ := rwd.voucherUC.GetVoucherCode(c, plValidator)

	if voucherCode != nil {
		rewards, err := rwd.voucherUC.VoucherValidate(c, plValidator)

		if err != nil {
			requestLogger.Debug(err)
			respErrors.SetTitle(models.ErrVoucherUnavailable.Error())

			return rwdInquiry, &respErrors
		}

		// get response reward
		rwdResp, _ := rwd.responseReward(c, rewards[0], plValidator)
		rwdInquiry.RefTrx = rwdResp.RefTrx
		rwdResp.RewardID = 0 // make rewardID nil
		rwdResp.RefTrx = ""

		if rwdResp != nil {
			rwdResponse = append(rwdResponse, *rwdResp)
		}

		rwdInquiry.Rewards = &rwdResponse

		// insert data to reward transaction
		_, err = rwd.createRewardTrx(c, *plValidator, rwdInquiry)

		if err != nil {
			requestLogger.Debug(err)
			respErrors.SetTitle(err.Error())

			return rwdInquiry, &respErrors
		}

		return rwdInquiry, &respErrors
	}

	// fresh or new reward trx start from here
	// check available campaign
	campaigns, err := rwd.campaignRepo.GetCampaignAvailable(c, *plValidator)

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)
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
		requestLogger.Debug(err)

		return rwdInquiry, &respErrors
	}

	// create array rewards
	rewards := rwd.putRewards(c, campaigns)

	for _, reward := range rewards {
		rewardLogger := logger.GetRequestLogger(c, reward.Validators)

		// validate promo code
		if err = rwd.validatePromoCode(*reward.Tags, reward, plValidator.PromoCode); err != nil {
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
		err = reward.Validators.Validate(plValidator)

		if err != nil {
			rewardLogger.Debug(err)
			respErrors.AddError(err.Error())

			continue
		}

		// get response reward
		rwdResp, err := rwd.responseReward(c, reward, plValidator)

		if err != nil {
			rewardLogger.Debug(err)
			respErrors.SetTitle(err.Error())

			return rwdInquiry, &respErrors
		}

		if rwdResp != nil {
			rwdResponse = append(rwdResponse, *rwdResp)
		}

		// update reward quota
		rwd.quotaUC.UpdateReduceQuota(c, reward.ID)

		// if not multi
		if plValidator.IsMulti == false {
			break
		}
	}

	// if no reward found
	if len(rwdResponse) == 0 {
		requestLogger.Debug(models.ErrMessageNoRewards)
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
		requestLogger.Debug(err)
		respErrors.SetTitle(err.Error())

		return rwdInquiry, &respErrors
	}

	// if not multi
	if plValidator.IsMulti == false {
		rwdInquiry.RefTrx = rwdResponse[0].RefTrx
		rwdResponse[0].RefTrx = ""
	}

	// check referral cant use referrer myself
	if plValidator.IsMulti == true &&(plValidator.CIF == plValidator.Referrer) {
		requestLogger.Debug(models.ErrSameCifReferrerAndReferral)
		requestLogger.Debug(err)
		respErrors.SetTitle(models.ErrSameCifReferrerAndReferral.Error())

		return rwdInquiry, &respErrors
	}

	return rwdInquiry, nil
}

func (rwd *rewardUseCase) responseReward(c echo.Context, reward models.Reward,
	plValidator *models.PayloadValidator) (*models.RewardResponse, error) {
	var rwdResp models.RewardResponse
	var voucherCode *models.VoucherCode
	var voucherName string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// get the rewards value/benefit
	rwdValue, _ := reward.Validators.GetRewardValue(plValidator)

	// check rewards voucher if any
	rwdVoucher, _ := reward.Validators.GetVoucherResult()

	if reward.RefID == "" {
		reward.RefID = randRefID(20)
	}

	if rwdVoucher != 0 {
		plVoucherBuy := &models.PayloadVoucherBuy{
			CIF:       plValidator.CIF,
			VoucherID: strconv.FormatInt(rwdVoucher, 10),
			RefID:     reward.RefID,
		}

		// check voucher code base on refId
		// if nil system will book a new voucher code
		voucherCodeValid, err := rwd.voucherCodeRepo.ValidateVoucherGive(c, plVoucherBuy)

		if voucherCodeValid != nil {
			voucherName = voucherCodeValid.Voucher.Name
		} else {
			// book a new voucher code
			voucherCode, err = rwd.voucherUC.VoucherGive(c, plVoucherBuy)

			if err != nil {
				requestLogger.Debug(models.ErrVoucherUnavailable)
				return nil, err
			}

			voucherName = voucherCode.Voucher.Name
		}

		rwdResp.VoucherName = voucherName

		rwdValue = 0 // if voucher reward is exist then reward value should be nil
	}

	// populate reward response
	rwdResp.Populate(reward, rwdValue, *plValidator)

	return &rwdResp, nil
}

func (rwd *rewardUseCase) Payment(c echo.Context, rwdPayment *models.RewardPayment) (models.RewardTrxResponse, error) {
	var responseData models.RewardTrxResponse
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	trimmedString := strings.Replace(rwdPayment.RefTrx, " ", "", -1)
	refIDs := strings.Split(trimmedString, ";")

	for _, refID := range refIDs {
		rwdPayment.RefTrx = refID

		// check available reward transaction based in ref_id
		rwdTrx, err := rwd.rwdTrxRepo.CheckRefID(c, rwdPayment.RefTrx)

		if err != nil {
			requestLogger.Debug(models.ErrRefTrxNotFound)

			return responseData, models.ErrRefTrxNotFound
		}

		// no ref_core equals to trx rejected
		if rwdPayment.RefCore == "" {
			// rejected
			// update voucher code
			rwd.voucherCodeRepo.UpdateVoucherCodeRejected(c, rwdPayment.RefTrx)

			// update reward trx
			rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, models.RewardTrxRejected)

			// update add reward quota
			rwd.quotaUC.UpdateAddQuota(c, *rwdTrx.RewardID)

			return responseData, nil
		}

		// trx status = succeeded or status = rejected or status = forceSucceeded then return error
		if *rwdTrx.Status != models.RewardTrxInquired && *rwdTrx.Status != models.RewardTrxTimeOut {
			responseData.StatusCode = rwdTrx.Status
			responseData.Status = rwdTrx.GetstatusRewardTrxText()

			return responseData, nil
		}

		// succeeded
		// update voucher code
		trxStatus := models.RewardTrxSucceeded
		rwd.voucherCodeRepo.UpdateVoucherCodeSucceeded(c, rwdPayment)

		if *rwdTrx.Status == models.RewardTrxTimeOut {
			// update reward trx timeout force to Succedeed
			trxStatus = models.RewardTrxTimeOutForceToSucceeded

			// update reduce reward quota
			rwd.quotaUC.UpdateReduceQuota(c, *rwdTrx.RewardID)
		}

		rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, trxStatus)

		// check if promoCode is a voucher code
		// if yes then it should redeem the voucher code itself
		pv := &models.PayloadValidator{PromoCode: rwdTrx.UsedPromoCode}
		voucherCode, _, _ := rwd.voucherUC.GetVoucherCode(c, pv)

		if voucherCode != nil {
			nowStr := time.Now().Format(models.DateTimeFormat)
			rwd.voucherCodeRepo.UpdateVoucherCodeRedeemed(c, nowStr, rwdPayment.CIF, rwdTrx.UsedPromoCode)
		}

		// check if a referral trx
		if rwdTrx.RequestData.IsReferral() {
			referralTrx := rwdTrx.GetReferralTrx()
			referralTrx.CIF = rwdPayment.CIF

			_ = rwd.referralTrxRepo.Create(c, referralTrx)
		}

		if rwdTrx.Reward.Type == nil {
			continue
		}

		// send sms notification only for voucher reward
		if *rwdTrx.Reward.Type == models.RewardTypeVoucher {
			go rwd.sendSmsVoucher(c, *rwdTrx)
		}
	}

	return responseData, nil
}

func (rwd *rewardUseCase) CheckTransaction(c echo.Context, rwdPayment *models.RewardPayment) (models.RewardTrxResponse, error) {
	var responseData models.RewardTrxResponse
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	trimmedString := strings.Replace(rwdPayment.RefTrx, " ", "", -1)
	refIDs := strings.Split(trimmedString, ";")

	for _, refID := range refIDs {
		rwdPayment.RefTrx = refID

		// check available reward transaction based in ref_id
		rewardtrx, err := rwd.rwdTrxRepo.CheckRefID(c, rwdPayment.RefTrx)

		if err != nil {
			requestLogger.Debug(models.ErrRefTrxNotFound)

			return responseData, models.ErrRefTrxNotFound
		}

		responseData.StatusCode = rewardtrx.Status
		responseData.Status = rewardtrx.GetstatusRewardTrxText()
	}

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

func (rwd *rewardUseCase) sendSmsVoucher(c echo.Context, rewardTrx models.RewardTrx) {
	var respBody map[string]interface{}
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	client := &http.Client{}
	voucherCode, err := rwd.voucherCodeRepo.GetVoucherCodeRefID(c, rewardTrx.RefID)

	if err != nil {
		requestLogger.Info(err)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
	}

	data := url.Values{}
	data.Set("message", fmt.Sprintf(models.VoucherSMSMessage, voucherCode.Voucher.Name,
		voucherCode.PromoCode, os.Getenv(`CS_NUMBER_1`), os.Getenv(`CS_NUMBER_2`)))
	data.Set("noHp", rewardTrx.RequestData.Phone)

	apiURL := os.Getenv(`PDS_API_HOST`) + os.Getenv(`SEND_SMS_PROMO_PATH`)

	if err != nil {
		requestLogger.Info(err)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv(`PDS_API_BASIC_USER`), os.Getenv(`PDS_API_BASIC_PASS`))
	logger.DataLog(c, data).Info("Start sending sms request to PDS API")
	response, err := client.Do(req)

	if err != nil || response == nil {
		requestLogger.Info(err)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
		logger.DataLog(c, respBody).Info("End sending sms request to PDS API")

		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		requestLogger.Info(err)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
		logger.DataLog(c, respBody).Info("End sending sms request to PDS API")

		return
	}

	err = json.Unmarshal(body, &respBody)

	if err != nil {
		requestLogger.Info(err)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
		logger.DataLog(c, respBody).Info("End sending sms request to PDS API")

		return
	}

	if respBody["status"] == "error" {
		requestLogger.Info(respBody)
		requestLogger.Info(models.DynamicErr(models.ErrSMSNotSent, rewardTrx.RefID))
		logger.DataLog(c, respBody).Info("End sending sms request to PDS API")

		return
	}

	logger.DataLog(c, respBody).Info("End sending sms request to PDS API")
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

func (rwd *rewardUseCase) validatePromoCode(tags []models.Tag, reward models.Reward, promoCode string) error {
	validPC := reward.PromoCode

	if *reward.IsPromoCode == models.IsPromoCodeFalse {
		return nil
	}

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

func (rwd *rewardUseCase) createRewardTrx(c echo.Context, plValidator models.PayloadValidator,
	rwdResponse models.RewardsInquiry) ([]*models.RewardTrx, error) {

	rewardTrx, err := rwd.rwdTrxRepo.Create(c, plValidator, rwdResponse)

	if err != nil {
		return nil, err
	}

	for _, trx := range rewardTrx {
		rwd.timeoutTrxJob(c, *trx)
	}

	return rewardTrx, nil
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

func (rwd *rewardUseCase) GetRewards(c echo.Context, rewardPayload *models.RewardsPayload) ([]models.Reward, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// counter data
	counter, err := rwd.rewardRepo.CountRewards(c, rewardPayload)
	if err != nil {
		requestLogger.Debug(models.ErrGetRewardCounter)

		return nil, "", models.ErrGetRewardCounter
	}

	// get data
	data, err := rwd.rewardRepo.GetRewards(c, rewardPayload)

	if err != nil {
		requestLogger.Debug(models.ErrGetReward)

		return nil, "", models.ErrGetReward
	}

	return data, strconv.FormatInt(counter, 10), err
}

func randRefID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = models.LetterBytes[rand.Int63()%int64(len(models.LetterBytes))]
	}

	return string(b)
}

func (rwd *rewardUseCase) validateReferralInq(c echo.Context, payload *models.PayloadValidator, respErrors *models.ResponseErrors) (bool, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if payload.Validators.CampaignCode != models.CampaignCodeReferral {
		return true, nil
	}

	modelsRefTrx := models.ReferralTrx{
		CIF:              payload.CIF,
		CifReferrer:      payload.Referrer,
		RefID:            payload.RefTrx,
		UsedReferralCode: payload.PromoCode,
		Type:             models.ReferralTrxTypeReferral,
		PhoneNumber:      payload.Phone,
	}

	// Value Reward CGC From GETENV
	rewardValue, err := strconv.Atoi(os.Getenv(`STAGE`))

	if err != nil {
		requestLogger.Debug(err)

		return false, err
	}

	// Limit Reward Counter Milestone From GETENV
	limitRewardCounter, err := strconv.Atoi(os.Getenv(`LIMIT_REWARD_COUNTER`))

	if err != nil {
		requestLogger.Debug(err)

		return false, err
	}

	countReftrx, err := rwd.referralTrxRepo.IsReferralTrxExist(c, modelsRefTrx)

	if err != nil {
		requestLogger.Debug(err)

		return false, err
	}

	if countReftrx > 0 {
		return false, models.ErrValidateGetReferral
	}

	totalGoldback, err := rwd.referralTrxRepo.GetTotalGoldbackReferrer(c, modelsRefTrx)

	if err != nil {
		requestLogger.Debug(err)

		return false, err
	}

	if int(totalGoldback) >= (rewardValue * limitRewardCounter) {
		return false, models.ErrValidateGetReferralMaxReward
	}

	return true, nil
}

func (rwd *rewardUseCase) GetRewardPromotions(c echo.Context, rplValidator models.RewardPromotionLists) ([]*models.RewardPromotions, *models.ResponseErrors, error) {
	var listPromotions []*models.RewardPromotions
	var err error
	logger := models.RequestLogger{}
	var respErrors models.ResponseErrors
	requestLogger := logger.GetRequestLogger(c, nil)

	listPromotions, err = rwd.rewardRepo.GetRewardPromotions(c, rplValidator)

	if err != nil {
		requestLogger.Debug(models.ErrGetRewardPromotions)

		respErrors.SetTitle(models.ErrGetRewardPromotions.Error())
		return nil, &respErrors, err
	}

	return listPromotions, &respErrors, nil
}
