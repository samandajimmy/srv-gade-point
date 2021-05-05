package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/pointhistories"
	"gade/srv-gade-point/vouchercodes"
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
)

var (
	floatType = reflect.TypeOf(float64(0))
)

type voucherUseCase struct {
	voucherRepo    vouchers.Repository
	vcRepo         vouchercodes.Repository
	campaignRepo   campaigns.Repository
	pHistoriesRepo pointhistories.Repository
}

// NewVoucherUseCase will create new an voucherUseCase object representation of vouchers.UseCase interface
func NewVoucherUseCase(vchrRepo vouchers.Repository, campgnRepo campaigns.Repository, pHistoriesRepo pointhistories.Repository, vcRepo vouchercodes.Repository) vouchers.UseCase {
	return &voucherUseCase{
		voucherRepo:    vchrRepo,
		campaignRepo:   campgnRepo,
		pHistoriesRepo: pHistoriesRepo,
		vcRepo:         vcRepo,
	}
}

func (vchr *voucherUseCase) CreateVoucher(c echo.Context, voucher *models.Voucher) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	err := vchr.voucherRepo.CreateVoucher(c, voucher)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherFailed)

		return models.ErrVoucherFailed
	}

	// create voucher codes that needed
	err = vchr.createVoucherCodes(c, voucher)

	if err != nil {
		return err
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
	apiURL := os.Getenv(`PDS_API_HOST`) + os.Getenv(`UPLOAD_IMAGE_URL`)

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
	response, err := http.Post(apiURL, contentType, bodyBuf)

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

	totalCount, err = vchr.voucherRepo.CountVouchers(c, false)

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

func (vchr *voucherUseCase) GetVouchers(c echo.Context) (models.ListVoucher, *models.ResponseErrors, error) {
	var listVoucher []*models.Voucher
	var responseVoucher models.ListVoucher
	var err error
	var totalCount int
	logger := models.RequestLogger{}
	var respErrors models.ResponseErrors
	requestLogger := logger.GetRequestLogger(c, nil)

	listVoucher, err = vchr.voucherRepo.GetVouchers(c)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		respErrors.SetTitle(models.ErrGetVouchers.Error())
		return responseVoucher, &respErrors, err
	}

	totalCount, err = vchr.voucherRepo.CountVouchers(c, true)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		respErrors.SetTitle(models.ErrGetVoucherCounter.Error())
		return responseVoucher, &respErrors, models.ErrGetVoucherCounter
	}

	if totalCount <= 0 {
		requestLogger.Debug(models.ErrVoucherUnavailable)

		respErrors.SetTitle(models.ErrVoucherUnavailable.Error())
		return responseVoucher, &respErrors, models.ErrVoucherUnavailable
	}

	responseVoucher.Vouchers = listVoucher
	responseVoucher.TotalCount = strconv.Itoa(totalCount)

	return responseVoucher, &respErrors, nil
}

func (vchr *voucherUseCase) GetHistoryVouchersUser(c echo.Context) (models.ListVoucher, error) {
	var err error
	var totalCount int
	var responseVoucher models.ListVoucher
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	vouchersUser, err := vchr.voucherRepo.GetHistoryVouchersUser(c)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return responseVoucher, err
	}

	totalCount, err = vchr.voucherRepo.CountHistoryVouchersUser(c, true)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return responseVoucher, models.ErrGetVoucherCounter
	}

	responseVoucher.Vouchers = vouchersUser
	responseVoucher.TotalCount = strconv.Itoa(totalCount)

	return responseVoucher, nil
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

func (vchr *voucherUseCase) GetVouchersUser(c echo.Context, payload map[string]interface{}) (models.ListVoucher, error) {
	var err error
	var totalCount int
	var responseVoucher models.ListVoucher
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	vouchersUser, err := vchr.voucherRepo.GetVouchersUser(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return responseVoucher, err
	}

	payload["status"] = "1" // to get the bought status
	payload["voucherID"] = ""

	totalCount, err = vchr.voucherRepo.CountVouchersUser(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCounter)

		return responseVoucher, models.ErrGetVoucherCounter
	}

	responseVoucher.Vouchers = vouchersUser
	responseVoucher.TotalCount = strconv.Itoa(totalCount)

	return responseVoucher, nil
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

	// get user current point
	userPoint, err := vchr.pHistoriesRepo.GetUserPoint(ech, payload.CIF)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPoint)

		return nil, models.ErrGetUserPoint
	}

	// validate voucher to buy
	err = validateBuy(voucherDetail.Point, int64(userPoint), voucherDetail.Available)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	voucherCode, err := vchr.voucherRepo.UpdatePromoCodeBought(ech, payload.VoucherID, payload.CIF)

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
		CIF:             payload.CIF,
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

func (vchr *voucherUseCase) VoucherGive(ech echo.Context, payload *models.PayloadVoucherBuy) (
	*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(ech, nil)
	now := time.Now()
	voucherDetail, err := vchr.voucherRepo.GetVoucher(ech, payload.VoucherID)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherUnavailable)

		return nil, models.ErrVoucherUnavailable
	}
	// create voucher codes after inquiry referral
	// target for referrer
	err = vchr.updateStockVoucher(ech, voucherDetail)

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

	// book voucher code if available
	voucherCode, err := vchr.voucherRepo.BookVoucherCode(ech, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUpdatePromoCodes)

		return nil, models.ErrUpdatePromoCodes
	}

	voucherCode.Voucher = &models.Voucher{
		ID:   voucherDetail.ID,
		Name: voucherDetail.Name,
	}

	return voucherCode, nil
}

func (vchr *voucherUseCase) GetVoucherCode(c echo.Context, pv *models.PayloadValidator) (*models.VoucherCode, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check voucher codes
	voucherCode, voucherID, err := vchr.voucherRepo.GetVoucherCode(c, pv)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherCodeUnavailable)

		return nil, "", models.ErrVoucherCodeUnavailable
	}

	return voucherCode, voucherID, nil
}

func (vchr *voucherUseCase) VoucherValidate(c echo.Context, pv *models.PayloadValidator) ([]models.Reward, error) {
	response := []models.Reward{}
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now, _ := time.Parse(models.DateTimeFormatMillisecond, pv.TransactionDate)

	// check voucher codes
	vc, voucherID, err := vchr.voucherRepo.GetVoucherCode(c, pv)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherCodeUnavailable)

		return response, models.ErrVoucherCodeUnavailable
	}

	// get voucher detail
	voucher, err := vchr.voucherRepo.GetVoucherAdmin(c, voucherID)

	if err != nil {
		requestLogger.Debug(err)

		return response, err
	}

	vStartDate, _ := time.Parse(time.RFC3339, voucher.StartDate)
	vEndDate, _ := time.Parse(time.RFC3339, voucher.EndDate)

	// check date expiry
	if vStartDate.After(now) {
		requestLogger.Debug(models.ErrVoucherNotStarted)

		return response, models.ErrVoucherNotStarted
	}

	if vEndDate.Before(now) {
		requestLogger.Debug(models.ErrVoucherExpired)

		return response, models.ErrVoucherExpired
	}

	// voucher validations
	err = voucher.Validators.Validate(pv)

	if err != nil {
		requestLogger.Debug(err)

		return response, err
	}

	vc.RefID = randRefID(20)
	_ = vchr.vcRepo.UpdateVoucherCodeInquired(c, *vc, *pv)

	response = append(response, models.Reward{})
	response[0].Validators = voucher.Validators
	response[0].Type = voucher.Type
	response[0].JournalAccount = voucher.JournalAccount
	response[0].ID = voucher.ID
	response[0].RefID = vc.RefID

	return response, nil
}

func (vchr *voucherUseCase) VoucherRedeem(c echo.Context, voucherRedeem *models.PayloadValidator) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	promoCode, err := vchr.voucherRepo.UpdatePromoCodeRedeemed(c, voucherRedeem.VoucherID, voucherRedeem.CIF, voucherRedeem.PromoCode)

	if err != nil {
		requestLogger.Debug(models.ErrRedeemVoucher)

		return nil, models.ErrRedeemVoucher
	}

	return promoCode, nil
}

func (vchr *voucherUseCase) updateStockVoucher(c echo.Context, voucherDetail *models.Voucher) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	oneStock := int32(1)
	if *voucherDetail.GeneratorType == models.VoucherTrxBased {
		voucherDetail.Stock = &oneStock
		/*
			CreateVoucherCode for [CGC]
		*/
		err := vchr.createVoucherCodes(c, voucherDetail)
		if err != nil {
			requestLogger.Debug(models.ErrVoucherGenearatePromoCodes)
			return err
		}
		/*
			Update Stock Voucher
		*/
		err = vchr.voucherRepo.UpdateVoucherStock(c, strconv.Itoa(int(voucherDetail.ID)))

		if err != nil {
			requestLogger.Debug(models.ErrVoucherGenearatePromoCodes)
			return err
		}

	}
	return nil
}

func (vchr *voucherUseCase) UpdateStatusBasedOnStartDate() error {
	err := vchr.voucherRepo.UpdateStatusBasedOnStartDate()

	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}

	return nil
}

func (vchr *voucherUseCase) createVoucherCodes(c echo.Context, voucher *models.Voucher) error {
	var promoCode []*models.VoucherCode
	now := time.Now()
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// generate the voucher codes based on the stock
	codes, err := generatePromoCode(voucher.Stock)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherGenearatePromoCodes)

		return models.ErrVoucherGenearatePromoCodes
	}

	// populate promo codes data
	if len(codes) > 0 {
		for i := 0; i < len(codes); i++ {
			pc := &models.VoucherCode{
				PromoCode: voucher.PrefixPromoCode + codes[i],
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

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}

	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func randRefID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = models.LetterBytes[rand.Int63()%int64(len(models.LetterBytes))]
	}

	return string(b)
}
