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
	ID             int64      `json:"id"`
	Name           string     `json:"name" validate:"required"`
	Description    string     `json:"description" validate:"required"`
	StartDate      string     `json:"startDate" validate:"required"`
	EndDate        string     `json:"endDate" validate:"required"`
	Point          int64      `json:"point" validate:"required"`
	JournalAccount string     `json:"journalAccount" validate:"required"`
	Value          float64    `json:"value" validate:"required"`
	ImageUrl       string     `json:"imageUrl" validate:"required"`
	Status         int8       `json:"status"`
	Validators     Validators `json:"validators"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type UpdateVoucher struct {
	Status int8 `json:"status"`
}

type PathVoucher struct {
	ImageUrl string `json:"imageUrl"`
}
