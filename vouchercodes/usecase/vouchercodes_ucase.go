package usecase

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchercodes"
	"gade/srv-gade-point/vouchers"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
)

var acceptedExtFiles = []string{".csv", ".json"}

type voucherCodeUseCase struct {
	voucherCodeRepo vouchercodes.VcRepository
	voucherRepo     vouchers.Repository
}

// NewVoucherCodeUseCase will create new an NewVoucherCodeUseCase object representation of vouchercode.UseCase interface
func NewVoucherCodeUseCase(vchrCodesRepo vouchercodes.VcRepository, vchrsRepo vouchers.Repository) vouchercodes.UseCase {
	return &voucherCodeUseCase{
		voucherCodeRepo: vchrCodesRepo,
		voucherRepo:     vchrsRepo,
	}
}

func (vchrCodeUs *voucherCodeUseCase) GetVoucherCodeHistory(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// count voucher by userId
	counter, err := vchrCodeUs.voucherCodeRepo.CountVoucherCode(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get voucher code history data
	data, err := vchrCodeUs.voucherCodeRepo.GetVoucherCodeHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherHistory)

		return nil, "", models.ErrGetVoucherHistory
	}

	return data, counter, err
}

func (vchrCodeUs *voucherCodeUseCase) GetVoucherCodes(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	// count voucher by voucherId
	counter, err := vchrCodeUs.voucherCodeRepo.CountVoucherCodeByVoucherID(c, payload)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get voucher codes
	data, err := vchrCodeUs.voucherCodeRepo.GetVoucherCodes(c, payload)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrGetVoucherCodes)

		return nil, "", models.ErrGetVoucherCodes
	}

	return data, counter, err
}

func (vchrCodeUs *voucherCodeUseCase) VoucherCodeRedeem(c echo.Context, voucherRedeem *models.PayloadValidator) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	promoCode, err := vchrCodeUs.voucherCodeRepo.UpdateVoucherCodeRedeemed(c, voucherRedeem.RedeemedDate, voucherRedeem.CIF, voucherRedeem.PromoCode)

	if err != nil {
		requestLogger.Debug(models.ErrRedeemVoucher)

		return nil, models.ErrRedeemVoucher
	}

	return promoCode, nil
}

func (vchrCodeUs *voucherCodeUseCase) GetBoughtVoucherCode(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// count voucher by criteria
	counter, err := vchrCodeUs.voucherCodeRepo.CountBoughtVoucherCode(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCodes)

		return nil, "", models.ErrGetVoucherCodes
	}

	// get voucher codes by criteria
	data, err := vchrCodeUs.voucherCodeRepo.GetBoughtVoucherCode(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCodes)

		return nil, "", models.ErrGetVoucherCodes
	}

	return data, counter, err
}

func (vchrCodeUs *voucherCodeUseCase) ImportVoucherCodes(c echo.Context, file *multipart.FileHeader, voucherID string) (string, error) {
	var vchrCodes []*models.VoucherCode
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check file
	if !contains(acceptedExtFiles, filepath.Ext(file.Filename)) {
		requestLogger.Debug(models.ErrAllowedExtVchrCodesImport)

		return "", models.ErrAllowedExtVchrCodesImport
	}

	// get voucher
	voucher, err := vchrCodeUs.voucherRepo.GetVoucherAdmin(c, voucherID)

	if err != nil {
		requestLogger.Debug(models.ErrGetVouchers)

		return "", models.ErrGetVouchers
	}

	// get voucher codes data
	if filepath.Ext(file.Filename) == ".json" {
		vchrCodes, err = convertJSONVoucherCodes(file, voucher)
	} else {
		vchrCodes, err = convertCSVVoucherCodes(file, voucher)
	}

	if err != nil {
		requestLogger.Debug(models.ErrMappingVchrCodesImport)

		return "", err
	}

	// store to db
	err = vchrCodeUs.voucherRepo.CreatePromoCode(c, vchrCodes, voucher.Synced)

	if err != nil {
		requestLogger.Debug(models.ErrVoucherStorePomoCodes)

		// Delete voucher when failed generate promo code
		err = vchrCodeUs.voucherRepo.DeleteVoucher(c, voucher.ID)

		if err != nil {
			requestLogger.Debug(models.ErrDeleteVoucher)
		}

		return "", models.ErrVoucherStorePomoCodes
	}

	return "", nil
}

func convertJSONVoucherCodes(file *multipart.FileHeader, voucher *models.Voucher) ([]*models.VoucherCode, error) {
	var vchrCodes []*models.VoucherCode
	var voucherCodes *models.VoucherCodes
	now := time.Now()

	// open file
	src, err := file.Open()

	if err != nil {

		return nil, err
	}

	defer src.Close()

	// read file
	byteValue, _ := ioutil.ReadAll(src)
	trimmedByte := bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))

	err = json.Unmarshal(trimmedByte, &voucherCodes)

	if err != nil {
		return nil, err
	}

	// data mapping
	for i := 0; i < len(voucherCodes.VoucherCodes); i++ {
		vchrCodes = append(vchrCodes, &models.VoucherCode{
			PromoCode: voucherCodes.VoucherCodes[i].PromoCode,
			Voucher:   voucher,
			CreatedAt: &now,
		})
	}

	return vchrCodes, err
}

func convertCSVVoucherCodes(file *multipart.FileHeader, voucher *models.Voucher) ([]*models.VoucherCode, error) {
	var voucherCodes []*models.VoucherCode
	now := time.Now()

	// open file
	csvFile, err := file.Open()

	if err != nil {

		return nil, err
	}

	defer csvFile.Close()

	// read file
	reader := csv.NewReader(csvFile)
	csvData, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	// data mapping
	for _, each := range csvData {
		voucherCode := &models.VoucherCode{
			PromoCode: each[0],
			Voucher:   voucher,
			CreatedAt: &now,
		}

		voucherCodes = append(voucherCodes, voucherCode)
	}

	return voucherCodes, err
}

func contains(strings []string, str string) bool {
	for _, n := range strings {
		if str == n {
			return true
		}
	}

	return false
}
