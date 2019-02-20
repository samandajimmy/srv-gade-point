package usecase

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
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

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (a *voucherUseCase) CreateVoucher(c context.Context, m *models.Voucher) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	err := a.voucherRepo.CreateVoucher(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *voucherUseCase) UpdateVoucher(c context.Context, id int64, updateVoucher *models.UpdateVoucher) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.voucherRepo.UpdateVoucher(ctx, id, updateVoucher)
	if err != nil {
		return err
	}

	return nil
}

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
	filePathPublic := os.Getenv(`VOUCHER_PATH`) + fileName

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

func (a *voucherUseCase) GetVoucher(c context.Context, name string, status string, startDate string, endDate string) (res interface{}, err error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listVoucher, err := a.voucherRepo.GetVoucher(ctx, name, status, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return listVoucher, nil
}
