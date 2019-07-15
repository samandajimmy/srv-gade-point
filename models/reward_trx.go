package models

import "time"

// RewardTrx is represent a reward_transactions model
type RewardTrx struct {
	ID            int64       `json:"id,omitempty"`
	Status        *int64      `json:"status,omitempty"`
	RefCore       string      `json:"refcore,omitempty"`
	RefID         string      `json:"refid,omitempty"`
	RewardID      string      `json:"rewardid,omitempty"`
	CIF           string      `json:"cif,omitempty"`
	UsedPromoCode string      `json:"usedPromoCode,omitempty"`
	InquiryDate   *time.Time  `json:"inquiryDate,omitempty"`
	SuccessedDate *time.Time  `json:"successedDate,omitempty"`
	RejectedDate  *time.Time  `json:"rejectedDate,omitempty"`
	TimeoutDate   *time.Time  `json:"timeoutDate,omitempty"`
	RequestData   RequestData `json:"requestData,omitempty"`
	CreatedAt     *time.Time  `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time  `json:"updatedAt,omitempty"`
}

// RequestData for JSONB
type RequestData map[string]interface{}
