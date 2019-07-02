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

	// ErrVoucherFailed to store create voucher failed error message
	ErrVoucherFailed = errors.New("Failed to create a voucher")

	// ErrVoucherGenearatePromoCodes to store generate promo codes error message
	ErrVoucherGenearatePromoCodes = errors.New("Something went wrong on generationg promotion codes")

	// ErrVoucherStorePomoCodes to store generate promo codes error message
	ErrVoucherStorePomoCodes = errors.New("Something went wrong on store promotion codes")

	// ErrDeleteVoucher to store delete voucher error message
	ErrDeleteVoucher = errors.New("Something went wrong when deleting a voucher")

	// ErrVoucherUpdateFailed to store update voucher failed error message
	ErrVoucherUpdateFailed = errors.New("Failed to update a voucher")

	// ErrOpenVoucherImg to store open file image
	ErrOpenVoucherImg = errors.New("Cannot open uploaded voucher image file")

	// ErrWriteVoucherImg to store open file image
	ErrWriteVoucherImg = errors.New("Cannot write uploaded voucher image file")

	// ErrCopyVoucherImg to store open file image
	ErrCopyVoucherImg = errors.New("Cannot copy uploaded voucher image file")

	// ErrGetVouchers to get vouchers error message
	ErrGetVouchers = errors.New("Something went wrong when trying to get vouchers")

	// ErrGetVoucherCodes to get vouchers error message
	ErrGetVoucherCodes = errors.New("Something went wrong when trying to get voucher codes")

	// ErrGetVoucherCounter to get voucher counter error message
	ErrGetVoucherCounter = errors.New("Something went wrong when trying to get voucher counter")

	// ErrRedeemVoucher to store redeem voucher error message
	ErrRedeemVoucher = errors.New("Voucher codes is not available to be redeemed")

	// ErrExceedBuyLimit to store exceeded buying limit voucher error message
	ErrExceedBuyLimit = errors.New("Exceeded buying limit for this voucher")

	// ErrVoucherNotStarted to store voucher not started error message
	ErrVoucherNotStarted = errors.New("Voucher has not started yet")

	// ErrVoucherUnavailable to store voucher unavailable error message
	ErrVoucherUnavailable = errors.New("Voucher Unavailable")

	// ErrBuyingVoucherExceeded to store exceeded error message
	ErrBuyingVoucherExceeded = errors.New("This user cant buy a voucher anymore today")

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

	// ErrValidationTrxAmt to store validation error message
	ErrValidationTrxAmt = errors.New("Your transaction amount is not enough to use this voucher")

	// ErrUsername to store username error message
	ErrUsername = errors.New("Username that you input is not valid")

	// ErrPassword to store password error message
	ErrPassword = errors.New("Password that you input is not valid")

	// ErrTokenExpired to store password error message
	ErrTokenExpired = errors.New("Your token has been expired")

	// ErrUsersNA to store users not available error message
	ErrUsersNA = errors.New("Users are not available")

	// ErrGetUsersPoint to store get users point error message
	ErrGetUsersPoint = errors.New("Something went wrong when trying to get users point")

	// ErrGetVoucherHistory to store get voucher history error message
	ErrGetVoucherHistory = errors.New("Something went wrong when trying to get voucher history")

	// ErrStartDateFormat to store a date format params error message
	ErrStartDateFormat = errors.New("Start date parameters is not meet the format")

	// ErrEndDateFormat to store a date format params error message
	ErrEndDateFormat = errors.New("End date parameters is not meet the format")

	// ErrAllowedExtVchrCodesImport to store allowed file ext error message
	ErrAllowedExtVchrCodesImport = errors.New("Import only allow csv and json file")

	// ErrMappingVchrCodesImport to store allowed file ext error message
	ErrMappingVchrCodesImport = errors.New("Something went wrong whent trying to map data from imported files")

	// ErrCreateMetric to store metric error message
	ErrCreateMetric = errors.New("Failed to create metric")

	// ErrUpdateMetric to store metric error message
	ErrUpdateMetric = errors.New("Failed to update metric")

	// ErrSendMetric to send metric error message
	ErrSendMetric = errors.New("Failed to send metric")
)
