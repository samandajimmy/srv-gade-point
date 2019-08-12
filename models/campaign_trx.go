package models

import "time"

// CampaignTrx is represent a campaign_transactions model
type CampaignTrx struct {
	ID              int64        `json:"id,omitempty"`
	CIF             string       `json:"cif,omitempty"`
	PointAmount     *float64     `json:"pointAmount,omitempty"`
	TransactionType string       `json:"transactionType,omitempty"`
	TransactionDate *time.Time   `json:"transactionDate,omitempty"`
	RefCore         string       `json:"refCore,omitempty"`
	RefID           string       `json:"refID,omitempty"`
	Campaign        *Campaign    `json:"campaign,omitempty"`
	VoucherCode     *VoucherCode `json:"voucherCode,omitempty"`
	UpdatedAt       *time.Time   `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time   `json:"createdAt,omitempty"`
}
