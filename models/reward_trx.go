package models

import "time"

var (
	// RewardTrxInquired to store reward inquired status
	RewardTrxInquired int64

	// RewardTrxSucceeded to store reward succeeded status
	RewardTrxSucceeded int64 = 1

	// RewardTrxRejected to store reward succeeded status
	RewardTrxRejected int64 = 2
)

// RewardTrx is represent a reward_transactions model
type RewardTrx struct {
	ID              int64      `json:"id,omitempty"`
	Status          *int64     `json:"status,omitempty"`
	RefCore         string     `json:"refCore,omitempty"`
	RefID           string     `json:"refId,omitempty"`
	RewardID        *int64     `json:"rewardId,omitempty"`
	CIF             string     `json:"cif,omitempty"`
	UsedPromoCode   string     `json:"usedPromoCode,omitempty"`
	TransactionDate *time.Time `json:"transactionDate,omitempty"`
	InquiredDate    *time.Time `json:"inquiredDate,omitempty"`
	SucceededDate   *time.Time `json:"succeededDate,omitempty"`
	RejectedDate    *time.Time `json:"rejectedDate,omitempty"`
	TimeoutDate     *time.Time `json:"timeoutDate,omitempty"`
	RequestData     string     `json:"requestData,omitempty"`
	ResponseData    string     `json:"responseData,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}
