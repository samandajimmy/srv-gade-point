package models

import "time"

// VoucherCode to store a voucher code data
type VoucherCode struct {
	ID           int64      `json:"id,omitempty"`
	PromoCode    string     `json:"promoCode,omitempty"`
	Status       *int8      `json:"status,omitempty"`
	UserID       string     `json:"userId,omitempty"`
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
