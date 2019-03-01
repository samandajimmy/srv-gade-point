package usecase

import (
	"context"
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
)

type voucherUseCase struct {
	voucherRepo    vouchers.Repository
	contextTimeout time.Duration
}

// NewVoucherUseCase will create new an voucherUseCase object representation of vouchers.UseCase interface
func NewVoucherUseCase(a vouchers.Repository, timeout time.Duration) vouchers.UseCase {
	return &voucherUseCase{
		voucherRepo:    a,
		contextTimeout: timeout,
	}
}

// create new voucher and generate promo code
func (a *voucherUseCase) CreateVoucher(c context.Context, m *models.Voucher) error {

	promoCode := make([]*models.PromoCode, 0)
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	err := a.voucherRepo.CreateVoucher(ctx, m)
	code, err := generatePromoCode(m.Stock)
	if err != nil {
		return err
	}

	for i := 0; i < len(code); i++ {
		ap := new(models.PromoCode)

		ap = &models.PromoCode{
			PromoCode: m.PrefixPromoCode + code[i],
			Status:    0,
			VoucherId: m.ID,
			CreatedAt: time.Now(),
		}
		promoCode = append(promoCode, ap)
	}

	err = a.voucherRepo.CreatePromoCode(ctx, promoCode)
	if err != nil {
		//Delete voucher when failed generate promo code
		err = a.voucherRepo.DeleteVoucher(ctx, m.ID)
		if err != nil {
			err = a.voucherRepo.DeleteVoucher(ctx, m.ID)
			return err
		}
		return err
	}

	return nil
}

// Update status voucher by id
func (a *voucherUseCase) UpdateVoucher(c context.Context, id int64, updateVoucher *models.UpdateVoucher) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.voucherRepo.UpdateVoucher(ctx, id, updateVoucher)
	if err != nil {
		return err
	}

	return nil
}

// Upload file image voucher
func (a *voucherUseCase) UploadVoucherImages(file *multipart.FileHeader) (string, error) {

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

//Get all voucher by param name, status, start date and end date
func (a *voucherUseCase) GetVouchers(c context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) (interface{}, string, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listVoucher, err := a.voucherRepo.GetVouchers(ctx, name, status, startDate, endDate, page, limit)
	if err != nil {
		return nil, "", err
	}

	totalCount, err := a.voucherRepo.CountVouchers(ctx, "")
	if err != nil {
		return nil, "", err
	}

	return listVoucher, strconv.Itoa(totalCount), nil
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
