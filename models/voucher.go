package models

import (
	"time"
)

type Validators struct {
	Channel         string `json:"channel" validate:"required"`
	Product         string `json:"product" validate:"required"`
	TransactionType string `json:"transactionType" validate:"required"`
	Unit            string `json:"unit" validate:"required"`
}

type Voucher struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name" validate:"required"`
	Description     string     `json:"description" validate:"required"`
	StartDate       string     `json:"startDate" validate:"required"`
	EndDate         string     `json:"endDate" validate:"required"`
	Point           int64      `json:"point" validate:"required"`
	JournalAccount  string     `json:"journalAccount" validate:"required"`
	Value           float64    `json:"value" validate:"required"`
	ImageUrl        string     `json:"imageUrl" validate:"required"`
	Status          int8       `json:"status"`
	Stock           int32      `json:"stock" validate:"required"`
	PrefixPromoCode string     `json:"prefixPromoCode" validate:"required"`
	Validators      Validators `json:"validators"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type UpdateVoucher struct {
	Status int8 `json:"status"`
}

type PathVoucher struct {
	ImageUrl string `json:"imageUrl"`
}

type VoucherDetail struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	StartDate   string  `json:"startDate" validate:"required"`
	EndDate     string  `json:"endDate" validate:"required"`
	Point       int64   `json:"point" validate:"required"`
	Value       float64 `json:"value" validate:"required"`
	ImageUrl    string  `json:"imageUrl" validate:"required"`
	Status      int8    `json:"status"`
	Stock       int32   `json:"stok"`
}

type PromoCode struct {
	ID           int64     `json:"id"`
	PromoCode    string    `json:"promoCode" validate:"required"`
	Status       int8      `json:"status"`
	UserId       string    `json:"userId" validate:"required"`
	VoucherId    int64     `json:"voucherId" validate:"required"`
	RedeemedDate time.Time `json:"redeemedDate"`
	UsedDate     time.Time `json:"UsedDate"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}
