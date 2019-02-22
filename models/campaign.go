package models

import (
	"time"
)

var (
	TransactionPointTypeDebet  = "D"
	TransactionPointTypeKredit = "K"
	CampaignTypePoint          = 0
)

type ValidatorsCampaign struct {
	Channel         string  `json:"channel" validate:"required"`
	Product         string  `json:"product" validate:"required"`
	TransactionType string  `json:"transactionType" validate:"required"`
	Unit            string  `json:"unit" validate:"required"`
	Multiplier      float64 `json:"multiplier" validate:"required"`
	Value           int64   `json:"value" validate:"required"`
	Formula         string  `json:"formula" validate:"required"`
}

type Campaign struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description" validate:"required"`
	StartDate   string             `json:"startDate" validate:"required"`
	EndDate     string             `json:"endDate" validate:"required"`
	Status      int8               `json:"status"`
	Type        int8               `json:"type"`
	Validators  ValidatorsCampaign `json:"validators"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt"`
}

type UpdateCampaign struct {
	Status int8 `json:"status"`
}

type GetCampaignValue struct {
	UserId            string  `json:"userId" validate:"required"`
	Channel           string  `json:"channel" validate:"required"`
	Product           string  `json:"product" validate:"required"`
	TransactionType   string  `json:"transactionType" validate:"required"`
	Unit              string  `json:"unit" validate:"required"`
	TransactionAmount float64 `json:"transactionAmount" validate:"required"`
}

type SaveTransactionPoint struct {
	ID              int64     `json:"id"`
	UserId          string    `json:"userId" validate:"required"`
	PointAmount     float64   `json:"pointAmount" validate:"required"`
	TransactionType string    `json:"transactionType" validate:"required"`
	TransactionDate time.Time `json:"transactionDate" validate:"required"`
	CampaingId      int64     `json:"campaignId" validate:"required"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedAt       time.Time `json:"createdAt"`
}

type UserPoint struct {
	UserPoint float64 `json:"userPoint"`
}
