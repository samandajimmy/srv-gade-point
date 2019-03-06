package models

import (
	"time"
)

// Validators is represent a validators inside voucher
type Validators struct {
	Channel            string `json:"channel,omitempty"`
	Product            string `json:"product,omitempty"`
	TransactionType    string `json:"transactionType,omitempty"`
	Unit               string `json:"unit,omitempty"`
	MinimalTransaction string `json:"minimalTransaction,omitempty"`
}

// Voucher is represent a vouchers model
type Voucher struct {
	ID              int64       `json:"id,omitempty"`
	Name            string      `json:"name,omitempty"`
	Description     string      `json:"description,omitempty"`
	StartDate       string      `json:"startDate,omitempty"`
	EndDate         string      `json:"endDate,omitempty"`
	Point           int64       `json:"point,omitempty"`
	JournalAccount  string      `json:"journalAccount,omitempty"`
	Value           *float64    `json:"value,omitempty"`
	ImageURL        string      `json:"imageUrl,omitempty"`
	Status          *int8       `json:"status,omitempty"`
	Stock           int32       `json:"stock,omitempty"`
	PrefixPromoCode string      `json:"prefixPromoCode,omitempty"`
	Amount          *int32      `json:"amount,omitempty"`
	Available       *int32      `json:"available,omitempty"`
	Bought          *int32      `json:"bought,omitempty"`
	Redeemed        *int32      `json:"redeemed,omitempty"`
	Expired         *int32      `json:"expired,omitempty"`
	Validators      *Validators `json:"validators,omitempty"`
	UpdatedAt       *time.Time  `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time  `json:"createdAt,omitempty"`
}

// UpdateVoucher to store payload for update status
type UpdateVoucher struct {
	Status int8 `json:"status,omitempty"`
}

// PathVoucher to store payload for upload voucher
type PathVoucher struct {
	ImageURL string `json:"imageUrl,omitempty"`
}

// PromoCode to store a promo codes data
type PromoCode struct {
	ID           int64      `json:"id,omitempty"`
	PromoCode    string     `json:"promoCode,omitempty"`
	Status       int8       `json:"status,omitempty"`
	UserID       string     `json:"userId,omitempty"`
	Voucher      *Voucher   `json:"voucher,omitempty"`
	RedeemedDate *time.Time `json:"redeemedDate,omitempty"`
	BoughtDate   *time.Time `json:"boughtDate,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
}

// PayloadVoucherBuy to store payload to buy a voucher
type PayloadVoucherBuy struct {
	VoucherID string `json:"voucherId,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

// PayloadValidateVoucher to store payload to validate a voucher
type PayloadValidateVoucher struct {
	PromoCode         string  `json:"promoCode,omitempty"`
	VoucherID         string  `json:"voucherId,omitempty"`
	UserID            string  `json:"userId,omitempty"`
	Channel           string  `json:"channel,omitempty"`
	Product           string  `json:"product,omitempty"`
	TransactionType   string  `json:"transactionType,omitempty"`
	Unit              string  `json:"unit,omitempty"`
	TransactionAmount float64 `json:"transactionAmount,omitempty"`
}
