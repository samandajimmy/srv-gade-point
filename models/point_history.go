package models

import "time"

var (
	// TransactionPointTypeDebet to store debet code
	TransactionPointTypeDebet = "D"
	// TransactionPointTypeKredit to store kredit code
	TransactionPointTypeKredit = "K"
)

// PointHistory is represent a campaign_transactions model
type PointHistory struct {
	ID              int64        `json:"id,omitempty"`
	CIF             string       `json:"cif,omitempty"`
	PointAmount     *float64     `json:"pointAmount,omitempty"`
	TransactionType string       `json:"transactionType,omitempty"`
	TransactionDate *time.Time   `json:"transactionDate,omitempty"`
	UsedFor         string       `json:"usedFor,omitempty"`
	RefCore         string       `json:"refCore,omitempty"`
	Status          *int64       `json:"status,omitempty"`
	Reward          *Reward      `json:"reward,omitempty"`
	VoucherCode     *VoucherCode `json:"voucherCode,omitempty"`
	UpdatedAt       *time.Time   `json:"updatedAt,omitempty"`
	CreatedAt       *time.Time   `json:"createdAt,omitempty"`
}

// UserPoint to store payload user point data
type UserPoint struct {
	UserPoint *float64 `json:"userPoint,omitempty"`
}
