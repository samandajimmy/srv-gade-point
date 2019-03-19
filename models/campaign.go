package models

import (
	"time"
)

var (
	// TransactionPointTypeDebet to store debet code
	TransactionPointTypeDebet = "D"
	// TransactionPointTypeKredit to store kredit code
	TransactionPointTypeKredit = "K"
)

// Campaign is represent a campaigns model
type Campaign struct {
	ID          int64      `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	StartDate   string     `json:"startDate,omitempty"`
	EndDate     string     `json:"endDate,omitempty"`
	Status      *int8      `json:"status,omitempty"`
	Type        *int8      `json:"type,omitempty"`
	Validators  *Validator `json:"validators,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

// CampaignTrx is represent a campaign_transactions model
type CampaignTrx struct {
	ID              int64      `json:"id,omitempty"`
	UserID          string     `json:"userId,omitempty"`
	PointAmount     *float64   `json:"pointAmount,omitempty"`
	TransactionType string     `json:"transactionType,omitempty"`
	TransactionDate *time.Time `json:"transactionDate,omitempty"`
	Campaign        *Campaign  `json:"campaign,omitempty"`
	PromoCode       *PromoCode `json:"promoCode,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
}

// UpdateCampaign to store payload update campaign
type UpdateCampaign struct {
	Status int8 `json:"status"`
}

// GetCampaignValue to store payload get campaign value
type GetCampaignValue struct {
	UserID            string  `json:"userId,omitempty"`
	Channel           string  `json:"channel,omitempty"`
	Product           string  `json:"product,omitempty"`
	TransactionType   string  `json:"transactionType,omitempty"`
	Unit              string  `json:"unit,omitempty"`
	TransactionAmount float64 `json:"transactionAmount,omitempty"`
	Source            string  `json:"source,omitempty"` // device name that user used
}

// UserPoint to store payload user point data
type UserPoint struct {
	UserPoint float64 `json:"userPoint,omitempty"`
}
