package models

import (
	"time"
)

var (
	// VoucherStockBased to store voucher generator type stock based
	VoucherStockBased int8
	// VoucherTrxBased to store voucher generator type trx based
	VoucherTrxBased int8 = 1
	// VoucherTypePoint to store reward type point
	VoucherTypePoint int64
	// VoucherTypeDirectDiscount to store reward type direct discount
	VoucherTypeDirectDiscount int64 = 1
	// VoucherTypePercentageDiscount to store reward type percentage discount
	VoucherTypePercentageDiscount int64 = 2
	// VoucherTypeGoldback to store reward type goldback
	VoucherTypeGoldback int64 = 3
	// VoucherTypeVoucher to store reward type voucher
	VoucherTypeVoucher int64 = 4
)

// ListVoucher is parent represent a vouchers model
type ListVoucher struct {
	Vouchers   []*Voucher `json:"listVoucher,omitempty"`
	TotalCount string     `json:"totalCount,omitempty"`
}

// Voucher is represent a vouchers model
type Voucher struct {
	ID                 int64      `json:"id,omitempty"`
	Name               string     `json:"name,omitempty"`
	Description        string     `json:"description,omitempty"`
	StartDate          string     `json:"startDate,omitempty"`
	EndDate            string     `json:"endDate"`
	Point              *int64     `json:"point,omitempty"`
	JournalAccount     string     `json:"journalAccount,omitempty"`
	Type               *int64     `json:"type,omitempty"`
	Value              *float64   `json:"value,omitempty"`
	ImageURL           string     `json:"imageUrl,omitempty"`
	Status             *int8      `json:"status,omitempty"`
	GeneratorType      *int8      `json:"generatorType,omitempty"`
	Stock              *int32     `json:"stock,omitempty"`
	PrefixPromoCode    string     `json:"prefixPromoCode,omitempty"`
	ProductCode        *string    `json:"productCode,omitempty"`
	TransactionType    *string    `json:"transactionType,omitempty"`
	MinLoanAmount      *string    `json:"minLoanAmount,omitempty"`
	PromoCode          *string    `json:"promoCode,omitempty"`
	Amount             *int32     `json:"amount,omitempty"`
	Available          *int32     `json:"available,omitempty"`
	Bought             *int32     `json:"bought,omitempty"`
	Redeemed           *int32     `json:"redeemed,omitempty"`
	Expired            *int32     `json:"expired,omitempty"`
	TermsAndConditions string     `json:"termsAndConditions,omitempty"`
	HowToUse           string     `json:"howToUse,omitempty"`
	Validators         *Validator `json:"validators,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
}

// UpdateVoucher to store payload for update status
type UpdateVoucher struct {
	Status    int64 `json:"status,omitempty"`
	VoucherID int64 `jsn:"voucherId, omitempty"`
}

// PathVoucher to store payload for upload voucher
type PathVoucher struct {
	ImageURL string `json:"imageUrl,omitempty"`
}

// PayloadVoucherBuy to store payload to buy a voucher
type PayloadVoucherBuy struct {
	VoucherID string   `json:"voucherId,omitempty"`
	CIF       string   `json:"CIF,omitempty"`
	RefID     string   `json:"refId,omitempty"`
	Voucher   *Voucher `json:"voucher,omitempty"`
}

// ResponseValidateVoucher to store response to validate a voucher
type ResponseValidateVoucher struct {
	Discount       *float64 `json:"discount,omitempty"`
	JournalAccount string   `json:"journalAccount,omitempty"`
}
