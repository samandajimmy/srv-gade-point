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
	ID                 int64       `json:"id,omitempty"`
	Name               string      `json:"name,omitempty"`
	Description        string      `json:"description,omitempty"`
	StartDate          string      `json:"startDate,omitempty"`
	EndDate            string      `json:"endDate,omitempty"`
	Point              *int64      `json:"point,omitempty"`
	JournalAccount     string      `json:"journalAccount,omitempty"`
	Value              *float64    `json:"value,omitempty"`
	ImageURL           string      `json:"imageUrl,omitempty"`
	Status             *int8       `json:"status,omitempty"`
	Stock              *int32      `json:"stock,omitempty"`
	PrefixPromoCode    string      `json:"prefixPromoCode,omitempty"`
	Amount             *int32      `json:"amount,omitempty"`
	Available          *int32      `json:"available,omitempty"`
	Bought             *int32      `json:"bought,omitempty"`
	Redeemed           *int32      `json:"redeemed,omitempty"`
	Expired            *int32      `json:"expired,omitempty"`
	TermsAndConditions string      `json:"termsAndConditions,omitempty"`
	HowToUse           string      `json:"howToUse,omitempty"`
	LimitPerUser       *int        `json:"limitPerUser,omitempty"`
	Validators         *Validators `json:"validators,omitempty"`
	UpdatedAt          *time.Time  `json:"updatedAt,omitempty"`
	CreatedAt          *time.Time  `json:"createdAt,omitempty"`
}

// UpdateVoucher to store payload for update status
type UpdateVoucher struct {
	Status int8 `json:"status,omitempty"`
}

// PathVoucher to store payload for upload voucher
type PathVoucher struct {
	ImageURL string `json:"imageUrl,omitempty"`
}

// PayloadVoucherBuy to store payload to buy a voucher
type PayloadVoucherBuy struct {
	VoucherID string `json:"voucherId,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

// PayloadValidateVoucher to store payload to validate a voucher
type PayloadValidateVoucher struct {
	PromoCode         string      `json:"promoCode,omitempty"`
	VoucherID         string      `json:"voucherId,omitempty"`
	UserID            string      `json:"userId,omitempty"`
	TransactionAmount float64     `json:"transactionAmount,omitempty"`
	Validators        *Validators `json:"validators,omitempty"`
}

// ResponseValidateVoucher to store response to validate a voucher
type ResponseValidateVoucher struct {
	Discount       *float64 `json:"discount,omitempty"`
	JournalAccount string   `json:"journalAccount,omitempty"`
}
