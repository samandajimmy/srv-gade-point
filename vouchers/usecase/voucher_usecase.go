package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lengthCode  = 5
	timeFormat  = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

var (
	statusVoucher = []string{"0", "1"} // voucher status
	floatType     = reflect.TypeOf(float64(0))
)

type voucherUseCase struct {
	voucherRepo    vouchers.Repository
	campaignRepo   campaigns.Repository
	contextTimeout time.Duration
}

// NewVoucherUseCase will create new an voucherUseCase object representation of vouchers.UseCase interface
func NewVoucherUseCase(vchrRepo vouchers.Repository, campgnRepo campaigns.Repository, timeout time.Duration) vouchers.UseCase {
	return &voucherUseCase{
		voucherRepo:    vchrRepo,
		campaignRepo:   campgnRepo,
		contextTimeout: timeout,
	}
}

func (vchr *voucherUseCase) CreateVoucher(c echo.Context, voucher *models.Voucher) error {
	now := time.Now()
	promoCode := make([]*models.VoucherCode, 0)
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	err := vchr.voucherRepo.CreateVoucher(c, voucher)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherFailed)

		return models.ErrVoucherFailed
	}

	code, err := generatePromoCode(voucher.Stock)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherGenearatePromoCodes)

		return models.ErrVoucherGenearatePromoCodes
	}

	// populate promo codes data
	if len(code) > 0 {
		for i := 0; i < len(code); i++ {
			pc := &models.VoucherCode{
				PromoCode: voucher.PrefixPromoCode + code[i],
				Voucher:   voucher,
				CreatedAt: &now,
			}
			promoCode = append(promoCode, pc)
		}
	}

	// store promo codes data
	err = vchr.voucherRepo.CreatePromoCode(c, promoCode)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherStorePomoCodes)

		// Delete voucher when failed generate promo code
		err = vchr.voucherRepo.DeleteVoucher(c, voucher.ID)

		if err != nil {
			requestLogger.Debug(models.ErrDeleteVoucher)
		}

		return models.ErrVoucherStorePomoCodes
	}

	return nil
}

func (vchr *voucherUseCase) UpdateVoucher(c echo.Context, id int64, updateVoucher *models.UpdateVoucher) error {
	now := time.Now()
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	voucherDetail, err := vchr.voucherRepo.GetVoucher(c, strconv.FormatInt(id, 10))

	if err != nil {
		requestLogger.Debug(models.ErrVoucherUnavailable)

		return models.ErrVoucherUnavailable
	}

	vEndDate, _ := time.Parse(time.RFC3339, voucherDetail.EndDate)

	if vEndDate.Before(now.Add(time.Hour * -24)) {
		requestLogger.Debug(models.ErrVoucherExpired)

		return models.ErrVoucherExpired
	}

	err = vchr.voucherRepo.UpdateVoucher(c, id, updateVoucher)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherUpdateFailed)

		return models.ErrVoucherUpdateFailed
	}

	return nil
}

func (vchr *voucherUseCase) UploadVoucherImages(c echo.Context, file *multipart.FileHeader) (string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	src, err := file.Open()

	if err != nil {
		requestLogger.Debug(err)

		return "", models.ErrOpenVoucherImg
	}

	defer src.Close()

	// upload image to pds server
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", file.Filename)

	if err != nil {
		requestLogger.Debug(err)

		return "", models.ErrWriteVoucherImg
	}

	_, err = io.Copy(fileWriter, src)
	if err != nil {
		requestLogger.Debug(err)

		return "", models.ErrCopyVoucherImg
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	requestLogger.Debug("Do a post request to pds api.")
	response, err := http.Post(os.Getenv(`UPLOAD_IMAGE_URL`), contentType, bodyBuf)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}
	value := gjson.Get(string(responseBody), "status")

	if value.String() != "success" {
		value = gjson.Get(string(responseBody), "message")
		requestLogger.Debug(errors.New(value.String()))

		return "", errors.New(value.String())

	}

	value = gjson.Get(string(responseBody), "data.filename")
	filePathPublic := value.String()

	return filePathPublic, nil
}

func (vchr *voucherUseCase) GetVouchersAdmin(c echo.Context, payload map[string]interface{}) ([]*models.Voucher, string, error) {
	var listVoucher []*models.Voucher
	var err error
	var totalCount int
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	_ = vchr.voucherRepo.UpdateExpiryDate(c)

	listVoucher, err = vchr.voucherRepo.GetVouchersAdmin(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return nil, "", models.ErrGetVouchers
	}

	totalCount, err = vchr.voucherRepo.CountVouchers(c, payload, false)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return nil, "", models.ErrGetVoucherCounter
	}

	if totalCount <= 0 {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return listVoucher, "", models.ErrGetVoucherCounter
	}

	return listVoucher, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) GetVoucherAdmin(c echo.Context, voucherID string) (*models.Voucher, error) {
	var voucherDetail *models.Voucher
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	voucherDetail, err := vchr.voucherRepo.GetVoucherAdmin(c, voucherID)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return nil, models.ErrGetVouchers
	}

	return voucherDetail, nil
}

func (vchr *voucherUseCase) GetVouchers(c echo.Context, payload map[string]interface{}) ([]*models.Voucher, string, error) {
	var listVoucher []*models.Voucher
	var err error
	var totalCount int
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	listVoucher, err = vchr.voucherRepo.GetVouchers(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return nil, "", err
	}

	totalCount, err = vchr.voucherRepo.CountVouchers(c, payload, false)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return nil, "", models.ErrGetVoucherCounter
	}

	if totalCount <= 0 {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return listVoucher, "", models.ErrGetVoucherCounter
	}

	return listVoucher, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) GetVoucher(c echo.Context, voucherID string) (*models.Voucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	voucherDetail, err := vchr.voucherRepo.GetVoucher(c, voucherID)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return nil, err
	}

	return voucherDetail, nil
}

func (vchr *voucherUseCase) GetVouchersUser(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	var err error
	var totalCount int
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	vouchersUser, err := vchr.voucherRepo.GetVouchersUser(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return nil, "", err
	}

	payload["status"] = "1" // to get the bought status
	payload["voucherID"] = ""
	totalCount, err = vchr.voucherRepo.CountPromoCode(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return nil, "", models.ErrGetVoucherCounter
	}

	if totalCount <= 0 {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return vouchersUser, "", models.ErrGetVoucherCounter
	}

	return vouchersUser, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) VoucherBuy(ech echo.Context, payload *models.PayloadVoucherBuy) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(ech, nil)
	now := time.Now()
	voucherDetail, err := vchr.voucherRepo.GetVoucher(ech, payload.VoucherID)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherUnavailable)

		return nil, models.ErrVoucherUnavailable
	}

	// check date expiry
	vStartDate, _ := time.Parse(time.RFC3339, voucherDetail.StartDate)
	vEndDate, _ := time.Parse(time.RFC3339, voucherDetail.EndDate)

	if vStartDate.After(now) {
		requestLogger.Debug(models.ErrVoucherExpired)

		return nil, models.ErrVoucherNotStarted
	}

	if vEndDate.Before(now) {
		requestLogger.Debug(models.ErrVoucherExpired)

		return nil, models.ErrVoucherExpired
	}

	// check voucher limit per user
	payloadPC := map[string]interface{}{
		"userID":    payload.UserID,
		"voucherID": payload.VoucherID,
		"status":    "", // to get whatever status of promo codes
	}

	userVchrCodesCount, err := vchr.voucherRepo.CountPromoCode(ech, payloadPC)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return nil, models.ErrGetVoucherCounter
	}

	if voucherDetail.LimitPerUser != nil &&
		*voucherDetail.LimitPerUser > 0 &&
		*voucherDetail.LimitPerUser <= userVchrCodesCount {
		requestLogger.Debug(models.ErrExceedBuyLimit)

		return nil, models.ErrExceedBuyLimit
	}

	// get user current point
	userPoint, err := vchr.campaignRepo.GetUserPoint(ech, payload.UserID)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPoint)

		return nil, models.ErrGetUserPoint
	}

	if userPoint == 0 {
		requestLogger.Debug(models.ErrUserPointNA)

		return nil, models.ErrUserPointNA
	}

	// validate voucher to buy
	err = validateBuy(voucherDetail.Point, int64(userPoint), voucherDetail.Available)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	voucherAmount, err := vchr.voucherRepo.CountBoughtVoucher(ech, payload.VoucherID, payload.UserID)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	// check limit per day
	err = validateLimitPurchase(*voucherDetail.DayPurchaseLimit, voucherAmount)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	voucherCode, err := vchr.voucherRepo.UpdatePromoCodeBought(ech, payload.VoucherID, payload.UserID)

	if err != nil {
		requestLogger.Debug(models.ErrUpdatePromoCodes)

		return nil, models.ErrUpdatePromoCodes
	}

	// Parse interface to float
	parseFloat, err := getFloat(voucherDetail.Point)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrVoucherPoint
	}

	if math.IsInf(parseFloat, 0) {
		requestLogger.Debug("the result of formula is infinity and beyond")

		return nil, models.ErrVoucherPoint
	}

	pointAmount := math.Floor(parseFloat)

	campaignTrx := &models.CampaignTrx{
		UserID:          payload.UserID,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeKredit,
		TransactionDate: &now,
		VoucherCode:     voucherCode,
		CreatedAt:       &now,
	}

	err = vchr.campaignRepo.SavePoint(ech, campaignTrx)

	if err != nil {
		requestLogger.Debug(models.ErrStoreCampaignTrx)

		return nil, models.ErrStoreCampaignTrx
	}

	voucherCode.Voucher = &models.Voucher{
		ID:   voucherDetail.ID,
		Name: voucherDetail.Name,
	}

	return voucherCode, nil
}

func (vchr *voucherUseCase) VoucherValidate(c echo.Context, validateVoucher *models.PayloadValidator) (*models.ResponseValidateVoucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()

	// check voucher codes
	_, voucherID, err := vchr.voucherRepo.GetVoucherCode(c, validateVoucher.PromoCode, validateVoucher.UserID)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherCodeUnavailable)

		return nil, models.ErrVoucherCodeUnavailable
	}

	// get voucher detail
	voucher, err := vchr.voucherRepo.GetVoucherAdmin(c, voucherID)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	vStartDate, _ := time.Parse(time.RFC3339, voucher.StartDate)
	vEndDate, _ := time.Parse(time.RFC3339, voucher.EndDate)

	// check date expiry
	if vStartDate.After(now) {
		requestLogger.Debug(models.ErrVoucherNotStarted)

		return nil, models.ErrVoucherNotStarted
	}

	if vEndDate.Before(now) {
		requestLogger.Debug(models.ErrVoucherExpired)

		return nil, models.ErrVoucherExpired
	}

	// voucher validations
	err = voucher.Validators.Validate(validateVoucher)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	responseValid := &models.ResponseValidateVoucher{
		Discount:       voucher.Validators.Value,
		JournalAccount: voucher.JournalAccount,
	}

	return responseValid, nil
}

func (vchr *voucherUseCase) VoucherRedeem(c echo.Context, voucherRedeem *models.PayloadValidator) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	promoCode, err := vchr.voucherRepo.UpdatePromoCodeRedeemed(c, voucherRedeem.VoucherID, voucherRedeem.UserID, voucherRedeem.PromoCode)

	if err != nil {
		requestLogger.Debug(models.ErrRedeemVoucher)

		return nil, models.ErrRedeemVoucher
	}

	return promoCode, nil
}

func (vchr *voucherUseCase) BadaiEmasGift(c echo.Context, plValidator *models.PayloadValidator) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// get badai emas voucher
	vouchers, err := vchr.voucherRepo.GetBadaiEmasVoucher(c)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrGetVouchers
	}

	// validate voucher by loan amount
	validVouchers := []*models.Voucher{}

	for _, voucher := range vouchers {
		err = voucher.Validators.Validate(plValidator)

		if err == nil {
			validVouchers = append(validVouchers, voucher)
		}
	}

	if len(validVouchers) < 1 {
		// no valid voucher available
		requestLogger.Debug(err)

		return nil, models.ErrVoucherUnavailable
	}

	// get latest voucher
	latestVoucher := validVouchers[0]

	// give the voucher to userID
	payloadVoucherBuy := &models.PayloadVoucherBuy{
		VoucherID: strconv.FormatInt(latestVoucher.ID, 10),
		UserID:    plValidator.UserID,
	}

	voucherCode, errMsg := vchr.VoucherBuy(c, payloadVoucherBuy)

	if err != nil {
		requestLogger.Debug(errMsg)

		return nil, errMsg
	}

	return voucherCode, nil
}

func generatePromoCode(stock *int32) (code []string, err error) {
	var arr = make([]string, *stock)
	arrChecker := map[string]bool{}

	for i := range arr {
		arr[i] = randStringBytes(lengthCode, arrChecker)
	}

	return arr, nil
}

func randStringBytes(n int, arrChecker map[string]bool) string {
	var randString string
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	randString = string(b)

	if arrChecker[randString] {
		randString = string(randStringBytes(lengthCode, arrChecker))
	}

	arrChecker[randString] = true

	return randString
}

func validateBuy(voucherPoint *int64, userPoint int64, avaliable *int32) error {
	if *avaliable <= 0 {
		return models.ErrVoucherOutOfStock
	}

	if userPoint < *voucherPoint {
		return models.ErrPointDeficit
	}

	return nil
}

func validateLimitPurchase(limitPurchase int64, voucherAmount int64) error {
	if limitPurchase <= 0 {
		return nil
	}

	if limitPurchase <= voucherAmount {
		return models.ErrBuyingVoucherExceeded
	}

	return nil
}

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}

	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func (vchr *voucherUseCase) UpdateStatusBasedOnStartDate() error {
	err := vchr.voucherRepo.UpdateStatusBasedOnStartDate()

	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}

	return nil
}
