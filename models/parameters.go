package models

import (
	"math"
	"time"

	"github.com/labstack/echo"
)

var (
	// StatusSuccess to store a status response success
	StatusSuccess = "Success"

	// StatusError to store a status response error
	StatusError = "Error"

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

	// DateTimeFormatZone to store a date time with zone format of timestamp
	DateTimeFormatZone = "2006-01-02T15:04:05Z"

	// DateTimeFormatMillisecond to store a date time format of timestamp to millisecond
	DateTimeFormatMillisecond = "2006-01-02 15:04:05.000"

	// DateFormat to store a date format of timestamp
	DateFormat = "2006-01-02"

	// DateFormatRegex to store a regex of dd/mm/yyyy date format
	DateFormatRegex = "(^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$)"

	// BatchSizeVoucherCodes to store a max length of data that need to be inserted for
	BatchSizeVoucherCodes = 21845

	// LetterBytes a string to generate random ID
	LetterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// VoucherSMSMessage is a message for voucher sms notification
	VoucherSMSMessage = "Selamat, Anda mendapatkan Voucher %s dengan kode %s Info: %s, %s"

	// ErrSMSNotSent to store error sms not sent
	ErrSMSNotSent = "Ref Transaksi: %s, Gagal mengirim SMS"

	// CifRefCodeExisted to notice if a cif already had referral codes
	CifRefCodeExisted = "Cif %s sudah memiliki Kode Referral"

	// MessageDataFound to store a message response of data found
	MessageDataFound = "Data ditemukan"

	ApiLogMessage = "API request to %s with endpoint %s"
)

// EchoGroup to store routes group
type EchoGroup struct {
	Admin    *echo.Group
	API      *echo.Group
	Token    *echo.Group
	Referral *echo.Group
}

// NowUTC to get real current datetime but UTC format
func NowUTC() time.Time {
	return time.Now().UTC().Add(7 * time.Hour)
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
