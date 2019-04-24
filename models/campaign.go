package models

import (
	"time"
)

// Campaign is represent a campaigns model
type Campaign struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	StartDate   string              `json:"startDate,omitempty"`
	EndDate     string              `json:"endDate,omitempty"`
	Status      *int8               `json:"status,omitempty"`
	Type        *int8               `json:"type,omitempty"`
	Validators  *ValidatorsCampaign `json:"validators,omitempty"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
	CreatedAt   *time.Time          `json:"createdAt,omitempty"`
}

// ValidatorsCampaign is represent a validators inside campaign
type ValidatorsCampaign struct {
	Channel         string  `json:"channel,omitempty"`
	Product         string  `json:"product,omitempty"`
	TransactionType string  `json:"transactionType,omitempty"`
	Unit            string  `json:"unit,omitempty"`
	Multiplier      float64 `json:"multiplier,omitempty"`
	Value           int64   `json:"value,omitempty"`
	Formula         string  `json:"formula,omitempty"`
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
	ReffCore          string  `json:"reffCore,omitempty"`
}

// UserPoint to store payload user point data
type UserPoint struct {
	UserPoint float64 `json:"userPoint,omitempty"`
}
