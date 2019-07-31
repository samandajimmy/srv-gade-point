package models

import "github.com/labstack/echo"

var (
	// StatusSuccess to store a status response success
	StatusSuccess = "Berhasil"

	// StatusError to store a status response error
	StatusError = "Gagal"

	// MessageSaveSuccess to store a success message response of save
	MessageSaveSuccess = "Berhasil Disimpan"

	// MessageUpdateSuccess to store a success message response of update
	MessageUpdateSuccess = "Berhasil Diperbaharui"

	// MessageUploadSuccess to store a success message response of upload
	MessageUploadSuccess = "Berhasil Unggah"

	// MessageDataSuccess to store a success message response of data
	MessageDataSuccess = "Data Berhasil Dikirim"

	// MessagePointSuccess to store a success message response of data
	MessagePointSuccess = "Data Berhasil Dikirim"

	// MessageUpdateError to store an errpr message response 0f update
	MessageUpdateError = "Gagal Mempebaharui"

	// MessageUploadError to store en erro message response of upload
	MessageUploadError = "Gagal Unggah"

	// MessageValidationError to store an error message response of field validation
	MessageValidationError = "Gagal Validasi Kolom"

	// MessageDataNotFound to store a message response of data not found
	MessageDataNotFound = "Data Tidak Ditemukan"

	// MessageUnprocessableEntity to store a message response of unproccessable entity
	MessageUnprocessableEntity = "Entitas Tidak Dapat Diproses"

	// MessageTokenFailed to store a message response token failure
	MessageTokenFailed = "Gagal Membuat Token!"

	// MicroTimeFormat to store a time format of micro timestamp
	MicroTimeFormat = "20060102150405.000000"

	// DateTimeFormat to store a date time format of timestamp
	DateTimeFormat = "2006-01-02 15:04:05"

	// DateTimeFormatMillisecond to store a date time format of timestamp to millisecond
	DateTimeFormatMillisecond = "2006-01-02 15:04:05.000"

	// DateFormat to store a date format of timestamp
	DateFormat = "2006-01-02"

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
