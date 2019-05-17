package models

import "time"

// PromoCode to store a promo codes data
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
