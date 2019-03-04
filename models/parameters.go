package models

import (
	"time"
)

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
	TimeNow                = time.Now()
)
