package models

import "time"

var (
	// TransactionPointTypeDebet to store debet code
	TransactionPointTypeDebet = "D"
	// TransactionPointTypeKredit to store kredit code
	TransactionPointTypeKredit = "K"
)

// CampaignTrx is represent a campaign_transactions model
type CampaignTrx struct {
	ID              int64      `json:"id,omitempty"`
	UserID          string     `json:"userId,omitempty"`
	PointAmount     *float64   `json:"pointAmount,omitempty"`
	TransactionType string     `json:"transactionType,omitempty"`
	TransactionDate *time.Time `json:"transactionDate,omitempty"`
	ReffCore        string     `json:"reffCore,omitempty"`
	Campaign        *Campaign  `json:"campaign,omitempty"`
	PromoCode       *PromoCode `json:"promoCode,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
}
