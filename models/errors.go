package models

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Your requested Item is not found")
	ErrConflict            = errors.New("Your Item already exist")
	ErrBadParamInput       = errors.New("Given Param is not valid")
	ErrNoCampaign          = errors.New("No campaign avaliable")
	ErrPointNotEnough      = errors.New("Point is not enough")
	ErrVoucherExpired      = errors.New("Voucher has been expired")
	ErrVoucherNotAvailable = errors.New("Voucher not available")
)
