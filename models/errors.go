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

	// ErrCampaignFailed to store create campaign failed error message
	ErrCampaignFailed = errors.New("Failed to create a campaign")

	// ErrCampaignUpdateFailed to store update campaign failed error message
	ErrCampaignUpdateFailed = errors.New("Failed to update a campaign")

	// ErrNoCampaign to store campaign not found error message
	ErrNoCampaign = errors.New("No campaign avaliable")

	// ErrGetCampaign to get campaign error message
	ErrGetCampaign = errors.New("Something went wrong when trying to get campaign")

	// ErrGetCampaignCounter to get campaign counter error message
	ErrGetCampaignCounter = errors.New("Something went wrong when trying to get campaign counter")

	// ErrCalculateFormulaCampaign to get campaign counter error message
	ErrCalculateFormulaCampaign = errors.New("Something went wrong when trying to calculate campaign formula")

	// ErrStoreCampaignTrx to get campaign counter error message
	ErrStoreCampaignTrx = errors.New("Something went wrong when trying to store campaign transaction")

	// ErrGetUserPoint to get user point error message
	ErrGetUserPoint = errors.New("Something went wrong when trying to get user point")

	// ErrUserPointNA to get user point N/A error message
	ErrUserPointNA = errors.New("You dont have any points yet")

	// ErrGetUserPointHistory to get user point history error message
	ErrGetUserPointHistory = errors.New("Something went wrong when trying to get user point histrory")

	// ErrUserPointHistoryNA to get user point history N/A error message
	ErrUserPointHistoryNA = errors.New("You dont have any points history yet")

	// ErrCampaignExpired to store campaign expired error message
	ErrCampaignExpired = errors.New("Campaign has been expired")

	// ErrPointDeficit to store point deficit error message
	ErrPointDeficit = errors.New("You dont have enough point to buy this voucher")

	// ErrVoucherExpired to store voucher expired error message
	ErrVoucherExpired = errors.New("Voucher has been expired")

	// ErrVoucherNotStarted to store voucher not started error message
	ErrVoucherNotStarted = errors.New("Voucher has not started yet")

	// ErrVoucherUnavailable to store voucher unavailable error message
	ErrVoucherUnavailable = errors.New("Voucher Unavailable")

	// ErrVoucherOutOfStock to store voucher unavailable error message
	ErrVoucherOutOfStock = errors.New("Voucher is out of stock")

	// ErrUpdatePromoCodes to update promo codes error message
	ErrUpdatePromoCodes = errors.New("Something went wrong when trying to update voucher promo codes")

	// ErrVoucherPoint to get voucher point error message
	ErrVoucherPoint = errors.New("Something went wrong when trying to get voucher point")

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
