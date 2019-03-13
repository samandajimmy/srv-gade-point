package models

import "errors"

var (
	// ErrInternalServerError to store internal server error message
	ErrInternalServerError = errors.New("Internal Server Error")

	// ErrNotFound to store not found error message
	ErrNotFound = errors.New("Your requested Item is not found")

	// ErrConflict to store conflicted error message
	ErrConflict = errors.New("Your Item already exist")

	// ErrBadParamInput to store bad parameter error message
	ErrBadParamInput = errors.New("Given Param is not valid")

	// ErrNoCampaign to store campaign not found error message
	ErrNoCampaign = errors.New("No campaign avaliable")

	// ErrPointDeficit to store point deficit error message
	ErrPointDeficit = errors.New("Point deficit")

	// ErrVoucherExpired to store voucher expired error message
	ErrVoucherExpired = errors.New("Voucher has been expired")

	// ErrVoucherNotStarted to store voucher not started error message
	ErrVoucherNotStarted = errors.New("Voucher has not started yet")

	// ErrVoucherUnavailable to store voucher unavailable error message
	ErrVoucherUnavailable = errors.New("Voucher Unavailable")

	// ErrVoucherCodeUnavailable to store voucher unavailable error message
	ErrVoucherCodeUnavailable = errors.New("Voucher code unavailable")

	// ErrValidatorUnavailable to store validator unavailable error message
	ErrValidatorUnavailable = errors.New("Validator is unavailable")

	// ErrValidation to store validation error message
	ErrValidation = errors.New("Some of your inputs are not valid")

	// ErrUsername to store username error message
	ErrUsername = errors.New("Username that you input is not valid")

	// ErrPassword to store password error message
	ErrPassword = errors.New("Password that you input is not valid")

	// ErrTokenExpired to store password error message
	ErrTokenExpired = errors.New("Your token has been expired")
)

type Errors struct {
	// do nothing
}

func (err *Errors) ErrValidation() error {
	return ErrValidation
}
