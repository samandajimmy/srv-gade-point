package models

import (
	"time"
)

var (
	TransactionPointTypeDebet  = "D"
	TransactionPointTypeKredit = "K"
	CampaignTypePoint          = 0
)

// Campaign is represent a campaigns model
type Campaign struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	StartDate   string              `json:"startDate,omitempty"`
	EndDate     string              `json:"endDate,omitempty"`
	Status      int8                `json:"status,omitempty"`
	Type        int8                `json:"type,omitempty"`
	Validators  *ValidatorsCampaign `json:"validators,omitempty"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
	CreatedAt   *time.Time          `json:"createdAt,omitempty"`
}

// ValidatorsCampaign is represent a validators inside campaign
type ValidatorsCampaign struct {
	Channel         string  `json:"channel" validate:"required"`
	Product         string  `json:"product" validate:"required"`
	TransactionType string  `json:"transactionType" validate:"required"`
	Unit            string  `json:"unit" validate:"required"`
	Multiplier      float64 `json:"multiplier" validate:"required"`
	Value           int64   `json:"value" validate:"required"`
	Formula         string  `json:"formula" validate:"required"`
}

// CampaignTrx is represent a campaign_transactions model
type CampaignTrx struct {
	ID              int64      `json:"id"`
	UserID          string     `json:"userId" validate:"required"`
	PointAmount     float64    `json:"pointAmount" validate:"required"`
	TransactionType string     `json:"transactionType" validate:"required"`
	TransactionDate *time.Time `json:"transactionDate" validate:"required"`
	Campaign        *Campaign  `json:"campaign,omitempty"`
	Voucher         *Voucher   `json:"voucher,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
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

type UserPoint struct {
	UserPoint float64 `json:"userPoint"`
}
