package usecase

import (
	"encoding/json"
	"fmt"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/logger"
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

func (rwd *rewardUseCase) responseInquiry(c echo.Context, reward models.Reward,
	plValidator *models.PayloadValidator) (*models.RewardResponse, error) {
	var rwdResp models.RewardResponse
	var voucherCode *models.VoucherCode
	var voucherName string
	var err error

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
		voucherCodeValid, _ := rwd.voucherCodeRepo.ValidateVoucherGive(c, plVoucherBuy)

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

func (rwd *rewardUseCase) Payment(c echo.Context, rwdPayment *models.RewardPayment) ([]models.RewardTrxResponse, error) {
	var responseData []models.RewardTrxResponse
	var errorAppend []error
	var rwdTrx *models.RewardTrx
	var err error
	var zero int64
	trimmedString := strings.Replace(rwdPayment.RefTrx, " ", "", -1)
	refIDs := strings.Split(trimmedString, ";")

	if rwdPayment.RootRefTrx != "" {
		refIDs, err = rwd.rwdTrxRepo.CheckRootRefId(c, rwdPayment.RootRefTrx)
	}

	if err != nil {
		return responseData, models.ErrRefTrxNotFound
	}

	responseDataOriginal := models.RewardTrxResponse{}
	for _, refID := range refIDs {
		responseDataOriginal = models.RewardTrxResponse{}
		rwdPayment.RefTrx = refID

		// check available reward transaction based in ref_id
		rwdTrx, err = rwd.rwdTrxRepo.CheckRefID(c, rwdPayment.RefTrx)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			errorAppend = append(errorAppend, models.ErrRefTrxNotFound)
			continue
		}

		// make rewardID = 0
		if rwdTrx.RewardID == nil {
			rwdTrx.RewardID = &zero
		}

		// no ref_core equals to trx rejected
		if rwdPayment.RefCore == "" {
			// rejected
			// update voucher code
			_ = rwd.voucherCodeRepo.UpdateVoucherCodeRejected(c, rwdPayment.RefTrx)

			// update reward trx
			_ = rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, models.RewardTrxRejected)

			// update add reward quota
			_ = rwd.quotaUC.UpdateAddQuota(c, *rwdTrx.RewardID)

			continue
		}

		// trx status = succeeded or status = rejected or status = forceSucceeded then return error
		if *rwdTrx.Status != models.RewardTrxInquired && *rwdTrx.Status != models.RewardTrxTimeOut {
			responseDataOriginal.StatusCode = rwdTrx.Status
			responseDataOriginal.Status = rwdTrx.GetstatusRewardTrxText()

			responseData = append(responseData, responseDataOriginal)
			continue
		}

		// succeeded
		// update voucher code
		trxStatus := models.RewardTrxSucceeded
		_ = rwd.voucherCodeRepo.UpdateVoucherCodeSucceeded(c, rwdPayment)

		if *rwdTrx.Status == models.RewardTrxTimeOut {
			// update reward trx timeout force to Succedeed
			trxStatus = models.RewardTrxTimeOutForceToSucceeded

			// update reduce reward quota
			_ = rwd.quotaUC.UpdateReduceQuota(c, *rwdTrx.RewardID)
		}

		_ = rwd.rwdTrxRepo.UpdateRewardTrx(c, rwdPayment, trxStatus)

		// check if promoCode is a voucher code
		// if yes then it should redeem the voucher code itself
		pv := &models.PayloadValidator{PromoCode: rwdTrx.UsedPromoCode}
		voucherCode, _, _ := rwd.voucherUC.GetVoucherCode(c, pv, true)

		if voucherCode != nil {
			nowStr := time.Now().Format(models.DateTimeFormat)
			_, _ = rwd.voucherCodeRepo.UpdateVoucherCodeRedeemed(c, nowStr, rwdPayment.CIF, rwdTrx.UsedPromoCode)
		}

		// check if a referral trx
		rewards := *rwdTrx.ResponseData.Rewards
		if rewards[0].Reference == models.RefTargetReferrer {
			referralTrx := rwdTrx.GetReferralTrx()
			referralTrx.CIF = rwdPayment.CIF

			_ = rwd.referralTrxRepo.RPostReferralTrx(c, referralTrx)
		}

		if rwdTrx.Reward.Type == nil {
			continue
		}

		// send sms notification only for voucher reward
		if *rwdTrx.Reward.Type == models.RewardTypeVoucher {
			go rwd.sendSmsVoucher(c, *rwdTrx)
		}
	}

	if len(errorAppend) > 0 {
		return responseData, errorAppend[0]
	}

	return responseData, nil
}

func (rwd *rewardUseCase) CheckTransaction(c echo.Context, rwdPayment *models.RewardPayment) (models.RewardTrxResponse, error) {
	var responseData models.RewardTrxResponse
	trimmedString := strings.Replace(rwdPayment.RefTrx, " ", "", -1)
	refIDs := strings.Split(trimmedString, ";")

	for _, refID := range refIDs {
		rwdPayment.RefTrx = refID
		// check available reward transaction based in ref_id
		rewardtrx, err := rwd.rwdTrxRepo.CheckRefID(c, rwdPayment.RefTrx)

		if err != nil {

			return responseData, models.ErrRefTrxNotFound
		}

		responseData.Status = rewardtrx.GetstatusRewardTrxText()
	}

	return responseData, nil
}

func (rwd *rewardUseCase) RefreshTrx() {
	now := models.NowUTC()
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

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
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
			_, _ = rwd.rewardRepo.GetRewardTags(c, &reward)
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
