package models

import "github.com/labstack/echo"

var (
	StatusSuccess          = "Success"
	StatusError            = "Error"
	MassageSaveSuccess     = "Successfully Saved"
	MassageUpdateSuccess   = "Successfully Updated"
	MassageUploadSuccess   = "Successfully Upload"
	MassagePointSuccess    = "Data Successfully Sent"
	MassageUpdateError     = "Update Error"
	MassageUploadError     = "Upload Failed"
	MassageValidationError = "Field validation"
	MassageForbiddenError  = "Forbidden access"
	MessageDataNotFound    = "Data Not Found"
	MicroTimeFormat        = "20060102150405.000000"
)

// EchoGroup to store routes group
type EchoGroup struct {
	Admin *echo.Group
	API   *echo.Group
	Token *echo.Group
}
