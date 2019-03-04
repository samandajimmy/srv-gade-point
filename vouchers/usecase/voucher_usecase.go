package usecase

import (
	"context"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lengthCode  = 5
	timeFormat  = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

var (
	statusVoucher = []string{"0", "1"} // voucher status
	sources       = []string{"admin"}
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

// create new voucher and generate promo code
func (vchr *voucherUseCase) CreateVoucher(c context.Context, m *models.Voucher) error {

	promoCode := make([]*models.PromoCode, 0)
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)

	defer cancel()

	err := vchr.voucherRepo.CreateVoucher(ctx, m)
	code, err := generatePromoCode(m.Stock)
	if err != nil {
		return err
	}

	for i := 0; i < len(code); i++ {
		ap := new(models.PromoCode)

		ap = &models.PromoCode{
			PromoCode: m.PrefixPromoCode + code[i],
			Status:    0,
			Voucher:   m,
			CreatedAt: time.Now(),
		}
		promoCode = append(promoCode, ap)
	}

	err = vchr.voucherRepo.CreatePromoCode(ctx, promoCode)
	if err != nil {
		//Delete voucher when failed generate promo code
		err = vchr.voucherRepo.DeleteVoucher(ctx, m.ID)
		if err != nil {
			err = vchr.voucherRepo.DeleteVoucher(ctx, m.ID)
			return err
		}
		return err
	}

	return nil
}

// Update status voucher by id
func (vchr *voucherUseCase) UpdateVoucher(c context.Context, id int64, updateVoucher *models.UpdateVoucher) error {

	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	err := vchr.voucherRepo.UpdateVoucher(ctx, id, updateVoucher)
	if err != nil {
		return err
	}

	return nil
}

// Upload file image voucher
func (vchr *voucherUseCase) UploadVoucherImages(file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()

	ext := filepath.Ext(file.Filename)
	nsec := time.Now().UnixNano() // number of nanoseconds unix
	fileName := strconv.FormatInt(nsec, 10) + ext

	filePathUpload := os.Getenv(`VOUCHER_UPLOAD_PATH`) + fileName
	filePathPublic := os.Getenv(`VOUCHER_PATH`) + "/" + fileName

	dst, err := os.Create(filePathUpload)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePathPublic, nil
}

// Get all voucher by param name, status, start date and end date
func (vchr *voucherUseCase) GetVouchers(c context.Context, name string, status string, startDate string, endDate string, page int32, limit int32, source string) (interface{}, string, error) {
	var listVoucher interface{}
	var err error
	var totalCount int
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	if source != sources[0] {
		listVoucher, err = vchr.voucherRepo.GetVouchersExternal(ctx, name, startDate, endDate, page, limit)
		if err != nil {
			return nil, "", err
		}

		totalCount, err = vchr.voucherRepo.CountVouchers(ctx, statusVoucher[1], true)
		if err != nil {
			return nil, "", err
		}

	} else {
		listVoucher, err = vchr.voucherRepo.GetVouchers(ctx, name, status, startDate, endDate, page, limit)
		if err != nil {
			return nil, "", err
		}

		totalCount, err = vchr.voucherRepo.CountVouchers(ctx, status, false)
		if err != nil {
			return nil, "", err
		}
	}

	return listVoucher, strconv.Itoa(totalCount), nil
}

// Get detail voucher
func (vchr *voucherUseCase) GetVoucher(c context.Context, voucherId string, source string) (interface{}, error) {
	var voucherDetail interface{}
	var err error
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	if source != sources[0] {
		voucherDetail, err = vchr.voucherRepo.GetVoucherExternal(ctx, voucherId)
		if err != nil {
			return nil, err
		}
	} else {
		voucherDetail, err = vchr.voucherRepo.GetVoucher(ctx, voucherId)
		if err != nil {
			return nil, err
		}

	}

	return voucherDetail, nil
}

// Get vouchers user
func (vchr *voucherUseCase) GetVouchersUser(c context.Context, userId string, status string, page int32, limit int32, source string) ([]*models.VoucherUser, string, error) {
	var err error
	var totalCount int
	ctx, cancel := context.WithTimeout(c, vchr.contextTimeout)
	defer cancel()

	vouchersUser, err := vchr.voucherRepo.GetVouchersUser(ctx, userId, status, page, limit)
	if err != nil {
		return nil, "", err
	}

	totalCount, err = vchr.voucherRepo.CountPromoCode(ctx, status, userId)
	if err != nil {
		return nil, "", err
	}

	return vouchersUser, strconv.Itoa(totalCount), nil
}

// Buy voucher
func (vchr *voucherUseCase) CreateVoucherBuy(c context.Context, m *models.PayloadVoucherBuy) (*models.VoucherUser, error) {
	// var err error

	// c, cancel := context.WithTimeout(c, vchr.contextTimeout)
	// defer cancel()

	// voucherDetail, err := vchr.voucherRepo.GetVoucher(c, m.VoucherId)
	// if err != nil {
	// 	return nil, err
	// }

	// userPoint, err := vchr.campaignRepo.GetUserPoint(c, m.UserId)
	// if err != nil {
	// 	return nil, err
	// }

	// EndDate, err := time.Parse(timeFormat, voucherDetail.EndDate)
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = validateBuy(EndDate, voucherDetail.Point, int64(userPoint), voucherDetail.Available)
	// if err != nil {
	// 	return nil, err
	// }

	// promoCodeId, promoCode, boughtDate, err := vchr.voucherRepo.UpdatePromoCodeBought(c, m.VoucherId, m.UserId)
	// if err != nil {
	// 	return nil, err
	// }

	// saveTransactionPoint := &models.SaveTransactionPoint{
	// 	UserId:          m.UserId,
	// 	PointAmount:     float64(voucherDetail.Point),
	// 	TransactionType: models.TransactionPointTypeKredit,
	// 	TransactionDate: time.Now(),
	// 	CampaingId:      0,
	// 	PromoCodeId:     promoCodeId,
	// 	CreatedAt:       time.Now(),
	// }

	// err = vchr.campaignRepo.SavePoint(c, saveTransactionPoint)
	// if err != nil {
	// 	return nil, err
	// }

	// voucherUser := &models.VoucherUser{
	// 	PromoCode:   promoCode,
	// 	BoughtDate:  boughtDate,
	// 	Name:        voucherDetail.Name,
	// 	Description: voucherDetail.Description,
	// 	Value:       voucherDetail.Value,
	// 	StartDate:   voucherDetail.StartDate,
	// 	EndDate:     voucherDetail.EndDate,
	// 	ImageUrl:    voucherDetail.ImageUrl,
	// }

	return nil, nil
}

// Generate promo code by stock, prefix code and length character code from data voucher
func generatePromoCode(stock int32) (code []string, err error) {

	var arr = make([]string, stock)
	for i := range arr {
		arr[i] = randStringBytes(lengthCode)
	}

	return arr, nil
}

// Rand String from letter bytes constant
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// validate buy voucher
func validateBuy(endDate time.Time, voucherPoint int64, userPoint int64, avaliable int32) (bool, error) {

	if endDate.Format(timeFormat) < time.Now().Format(timeFormat) {
		return false, models.ErrVoucherExpired
	}
	if avaliable <= 0 {
		return false, models.ErrVoucherUnavailable
	}
	if userPoint < voucherPoint {
		return false, models.ErrPointDeficit
	}

	return true, nil
}
