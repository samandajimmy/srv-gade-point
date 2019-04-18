package usecase

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/iancoleman/strcase"
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

func (vchr *voucherUseCase) CreateVoucher(c context.Context, m *models.Voucher) error {
	now := time.Now()
	promoCode := make([]*models.PromoCode, 0)
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)

	defer cancel()

	err := vchr.voucherRepo.CreateVoucher(ctx, m)
	code, err := generatePromoCode(m.Stock)

	if err != nil {
		return err
	}

	if len(code) > 0 {
		for i := 0; i < len(code); i++ {
			pc := new(models.PromoCode)
			pc = &models.PromoCode{
				PromoCode: m.PrefixPromoCode + code[i],
				Voucher:   m,
				CreatedAt: &now,
			}
			promoCode = append(promoCode, pc)
		}
	}

	err = vchr.voucherRepo.CreatePromoCode(ctx, promoCode)

	if err != nil {
		//Delete voucher when failed generate promo code
		err = vchr.voucherRepo.DeleteVoucher(ctx, m.ID)
		return err
	}

	return nil
}

func (vchr *voucherUseCase) UpdateVoucher(c context.Context, id int64, updateVoucher *models.UpdateVoucher) error {
	var voucherDetail *models.Voucher
	now := time.Now()
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	voucherDetail, err := vchr.voucherRepo.GetVoucher(ctx, strconv.FormatInt(id, 10))

	if voucherDetail == nil {
		log.Error(models.ErrVoucherUnavailable)
		return models.ErrVoucherUnavailable
	}

	vEndDate, _ := time.Parse(time.RFC3339, voucherDetail.EndDate)

	if vEndDate.Before(now.Add(time.Hour * -24)) {
		log.Error(models.ErrVoucherExpired)
		return models.ErrVoucherExpired
	}

	err = vchr.voucherRepo.UpdateVoucher(ctx, id, updateVoucher)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (vchr *voucherUseCase) UploadVoucherImages(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()

	if err != nil {
		return "", err
	}

	defer src.Close()

	// upload image to pds server
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", file.Filename)

	if err != nil {
		log.Error(err)
		return "", err
	}

	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	response, err := http.Post(os.Getenv(`UPLOAD_IMAGE_URL`), contentType, bodyBuf)

	if err != nil {
		log.Error(err)

		return "", err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Error(err)

		return "", err
	}
	value := gjson.Get(string(responseBody), "status")

	if value.String() != "success" {
		value = gjson.Get(string(responseBody), "message")
		log.Error(errors.New(value.String()))

		return "", errors.New(value.String())

	}

	value = gjson.Get(string(responseBody), "data.filename")
	filePathPublic := value.String()

	return filePathPublic, nil
}

func (vchr *voucherUseCase) GetVouchersAdmin(c context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, string, error) {
	var listVoucher []*models.Voucher
	var err error
	var totalCount int
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	err = vchr.voucherRepo.UpdateExpiryDate(ctx)
	if err != nil {
		log.Warn(err)
	}

	listVoucher, err = vchr.voucherRepo.GetVouchersAdmin(ctx, name, status, startDate, endDate, page, limit)

	if err != nil {
		return nil, "", err
	}

	totalCount, err = vchr.voucherRepo.CountVouchers(ctx, name, status, startDate, endDate, false)

	if err != nil {
		return nil, "", err
	}

	return listVoucher, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) GetVoucherAdmin(c context.Context, voucherID string) (*models.Voucher, error) {
	var voucherDetail *models.Voucher
	var err error
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	voucherDetail, err = vchr.voucherRepo.GetVoucherAdmin(ctx, voucherID)

	if err != nil {
		return nil, err
	}

	return voucherDetail, nil
}

func (vchr *voucherUseCase) GetVouchers(c context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, string, error) {
	var listVoucher []*models.Voucher
	var err error
	var totalCount int

	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	listVoucher, err = vchr.voucherRepo.GetVouchers(ctx, name, startDate, endDate, page, limit)

	if err != nil {
		return nil, "", err
	}

	totalCount, err = vchr.voucherRepo.CountVouchers(ctx, name, statusVoucher[1], startDate, endDate, true)

	if err != nil {
		return nil, "", err
	}

	return listVoucher, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) GetVoucher(c context.Context, voucherID string) (*models.Voucher, error) {
	var voucherDetail *models.Voucher
	var err error
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	voucherDetail, err = vchr.voucherRepo.GetVoucher(ctx, voucherID)

	if err != nil {
		return nil, err
	}

	return voucherDetail, nil
}

func (vchr *voucherUseCase) GetVouchersUser(c context.Context, userID string, status string, page int, limit int) ([]models.PromoCode, string, error) {
	var err error
	var totalCount int
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	vouchersUser, err := vchr.voucherRepo.GetVouchersUser(ctx, userID, status, page, limit)

	if err != nil {
		return nil, "", err
	}

	totalCount, err = vchr.voucherRepo.CountPromoCode(ctx, status, userID)

	if err != nil {
		return nil, "", err
	}

	return vouchersUser, strconv.Itoa(totalCount), nil
}

func (vchr *voucherUseCase) VoucherBuy(c context.Context, m *models.PayloadVoucherBuy) (*models.PromoCode, error) {
	var err error
	now := time.Now()
	c, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	err = vchr.voucherRepo.VoucherCheckExpired(c, m.VoucherID)

	if err != nil {
		return nil, err
	}

	voucherDetail, err := vchr.voucherRepo.GetVoucher(c, m.VoucherID)

	if err != nil {
		return nil, err
	}

	userPoint, err := vchr.campaignRepo.GetUserPoint(c, m.UserID)

	if err != nil {
		return nil, err
	}

	err = validateBuy(voucherDetail.Point, int64(userPoint), voucherDetail.Available)

	if err != nil {
		return nil, err
	}

	promoCode, err := vchr.voucherRepo.UpdatePromoCodeBought(c, m.VoucherID, m.UserID)

	if err != nil {
		return nil, err
	}

	// Parse interface to float
	parseFloat, err := getFloat(voucherDetail.Point)
	pointAmount := math.Floor(parseFloat)

	campaignTrx := &models.CampaignTrx{
		UserID:          m.UserID,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeKredit,
		TransactionDate: &now,
		PromoCode:       promoCode,
		CreatedAt:       &now,
	}

	err = vchr.campaignRepo.SavePoint(c, campaignTrx)

	if err != nil {
		return nil, err
	}

	promoCode.Voucher = voucherDetail
	return promoCode, nil
}

func (vchr *voucherUseCase) VoucherValidate(c context.Context, validateVoucher *models.PayloadValidateVoucher) (*models.ResponseValidateVoucher, error) {
	var payloadValidator map[string]interface{}
	now := time.Now()
	c, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	// check voucher codes
	_, voucherID, err := vchr.voucherRepo.GetVoucherCode(c, validateVoucher.PromoCode, validateVoucher.UserID)

	if err != nil {
		log.Error(err)
		return nil, models.ErrVoucherCodeUnavailable
	}

	// get voucher detail
	voucher, err := vchr.voucherRepo.GetVoucherAdmin(c, voucherID)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	vStartDate, _ := time.Parse(time.RFC3339, voucher.StartDate)
	vEndDate, _ := time.Parse(time.RFC3339, voucher.EndDate)

	// check date expiry
	if vStartDate.After(now) {
		log.Error(models.ErrVoucherNotStarted)
		return nil, models.ErrVoucherNotStarted
	}

	if vEndDate.Before(now) {
		log.Error(models.ErrVoucherExpired)
		return nil, models.ErrVoucherExpired
	}

	// voucher validations
	validator := voucher.Validators

	if validator == nil {
		log.Error(models.ErrValidatorUnavailable)
		return nil, models.ErrValidatorUnavailable
	}

	vReflector := reflect.ValueOf(validator).Elem()
	tempJSON, _ := json.Marshal(validateVoucher.Validators)
	json.Unmarshal(tempJSON, &payloadValidator)

	for i := 0; i < vReflector.NumField(); i++ {
		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue := vReflector.Field(i).Interface()

		switch fieldName {
		case "channel", "product", "transactionType", "unit":
			if fieldValue != payloadValidator[fieldName] {
				log.Error(models.ErrValidation)
				return nil, models.ErrValidation
			}
		case "minimalTransaction":
			minTrx, _ := strconv.ParseFloat(fieldValue.(string), 64)

			if minTrx > validateVoucher.TransactionAmount {
				log.Error(models.ErrValidation)
				return nil, models.ErrValidation
			}
		}
	}

	responseValid := &models.ResponseValidateVoucher{
		Discount:       voucher.Value,
		JournalAccount: voucher.JournalAccount,
	}

	return responseValid, nil
}

func (vchr *voucherUseCase) VoucherRedeem(c context.Context, voucherRedeem *models.PayloadValidateVoucher) (*models.PromoCode, error) {
	var err error
	c, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()
	promoCode, err := vchr.voucherRepo.UpdatePromoCodeRedeemed(c, voucherRedeem.VoucherID, voucherRedeem.UserID)

	if err != nil {
		return nil, err
	}

	return promoCode, nil
}

func generatePromoCode(stock *int32) (code []string, err error) {
	var arr = make([]string, *stock)

	for i := range arr {
		arr[i] = randStringBytes(lengthCode)
	}

	return arr, nil
}

func randStringBytes(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func validateBuy(voucherPoint *int64, userPoint int64, avaliable *int32) error {
	if *avaliable <= 0 {
		return models.ErrVoucherUnavailable
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

func (vchr *voucherUseCase) UpdateStatusBasedOnStartDate() error {

	err := vchr.voucherRepo.UpdateStatusBasedOnStartDate()
	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}
	return nil
}
