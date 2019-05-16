package models

import "github.com/labstack/echo"

var (
	// StatusSuccess to store a status response success
	StatusSuccess = "Success"

	// StatusError to store a status response error
	StatusError = "Error"

	// MessageSaveSuccess to store a success message response of save
	MessageSaveSuccess = "Successfully Saved"

	// MessageUpdateSuccess to store a success message response of update
	MessageUpdateSuccess = "Successfully Updated"

	// MessageUploadSuccess to store a success message response of upload
	MessageUploadSuccess = "Successfully Upload"

	// MessageDataSuccess to store a success message response of data
	MessageDataSuccess = "Data Successfully Sent"

	// MessagePointSuccess to store a success message response of data
	MessagePointSuccess = "Data Successfully Sent"

	// MessageUpdateError to store an errpr message response 0f update
	MessageUpdateError = "Update Error"

	// MessageUploadError to store en erro message response of upload
	MessageUploadError = "Upload Failed"

	// MessageValidationError to store an error message response of field validation
	MessageValidationError = "Field validation error"

	// MessageDataNotFound to store a message response of data not found
	MessageDataNotFound = "Data Not Found"

	// MessageUnprocessableEntity to store a message response of unproccessable entity
	MessageUnprocessableEntity = "Unprocessable Entity"

	// MessageTokenFailed to store a message response token failure
	MessageTokenFailed = "Creating token failed!"

	// MicroTimeFormat to store a time format of micro timestamp
	MicroTimeFormat = "20060102150405.000000"

	// DateFormatRegex to store a regex of dd/mm/yyyy date format
	DateFormatRegex = "(^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$)"

	// BatchSizeVoucherCodes to store a max length of data that need to be inserted for
	BatchSizeVoucherCodes = 21845
)

// EchoGroup to store routes group
type EchoGroup struct {
	Admin *echo.Group
	API   *echo.Group
	Token *echo.Group
}
