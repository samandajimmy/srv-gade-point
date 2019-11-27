package models

import "time"

var (
	// VoucherCodeStatusAvailable to store voucher code available status
	VoucherCodeStatusAvailable int64

	// VoucherCodeStatusBought to store voucher code bought status
	VoucherCodeStatusBought int64 = 1

	// VoucherCodeStatusRedeemed to store voucher code redeemed status
	VoucherCodeStatusRedeemed int64 = 2

	// VoucherCodeStatusExpired to store voucher code expired status
	VoucherCodeStatusExpired int64 = 3

	// VoucherCodeStatusBooked to store voucher code booked status
	VoucherCodeStatusBooked int64 = 4

	// VoucherCodeStatusInquired to store voucher code inquired status
	VoucherCodeStatusInquired int64 = 5
)

// VoucherCode to store a voucher code data
type VoucherCode struct {
	ID           int64      `json:"id,omitempty"`
	PromoCode    string     `json:"promoCode,omitempty"`
	Status       *int8      `json:"status,omitempty"`
	UserID       string     `json:"userId,omitempty"`
	RefID        string     `json:"refId,omitempty"`
	Voucher      *Voucher   `json:"voucher,omitempty"`
	RedeemedDate *time.Time `json:"redeemedDate,omitempty"`
	BoughtDate   *time.Time `json:"boughtDate,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
}

// VoucherCodes to a voucher codes data
type VoucherCodes struct {
	VoucherCodes []VoucherCode `json:"voucherCodes"`
}
