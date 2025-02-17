package models

import "time"

var (
	// RewardTrxInquired to store reward inquired status
	RewardTrxInquired int64

	// RewardTrxSucceeded to store reward succeeded status
	RewardTrxSucceeded int64 = 1

	// RewardTrxRejected to store reward succeeded status
	RewardTrxRejected int64 = 2

	// RewardTrxTimeOut to store reward time out status
	RewardTrxTimeOut int64 = 3

	// RewardTrxTimeOutForceToSucceeded to store reward time out force to success status
	RewardTrxTimeOutForceToSucceeded int64 = 4
)

var statusRewardTrx = map[int64]string{
	RewardTrxInquired:                "inquiry",
	RewardTrxSucceeded:               "succeded",
	RewardTrxRejected:                "rejected",
	RewardTrxTimeOut:                 "expired",
	RewardTrxTimeOutForceToSucceeded: "forced-succeded",
}

// RewardTrx is represent a reward_transactions model
type RewardTrx struct {
	ID              int64             `json:"id,omitempty"`
	Status          *int64            `json:"status,omitempty"`
	RefCore         string            `json:"refCore,omitempty"`
	RefID           string            `json:"refId,omitempty"`
	RewardID        *int64            `json:"rewardId,omitempty"`
	CIF             string            `json:"cif,omitempty"`
	UsedPromoCode   string            `json:"usedPromoCode,omitempty"`
	TransactionDate *time.Time        `json:"transactionDate,omitempty"`
	InquiredDate    *time.Time        `json:"inquiredDate,omitempty"`
	SucceededDate   *time.Time        `json:"succeededDate,omitempty"`
	RejectedDate    *time.Time        `json:"rejectedDate,omitempty"`
	TimeoutDate     *time.Time        `json:"timeoutDate,omitempty"`
	RequestData     *PayloadValidator `json:"requestData,omitempty"`
	ResponseData    *RewardsInquiry   `json:"responseData,omitempty"`
	CreatedAt       *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time        `json:"updatedAt,omitempty"`
	Reward          *Reward           `json:"reward,omitempty"`
}

// RewardTrxResponse is represent a reward transaction response model
type RewardTrxResponse struct {
	StatusCode *int64 `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
}

// GetstatusRewardTrxText to get text of status reward transaction
func (rwdTrx RewardTrx) GetstatusRewardTrxText() string {
	return statusRewardTrx[*rwdTrx.Status]
}

// GetReferralTrx to get all referral transaction
func (rwdTrx RewardTrx) GetReferralTrx() ReferralTrx {
	refTrx := ReferralTrx{}
	reward := *rwdTrx.ResponseData.Rewards

	refTrx.CIF = reward[0].CIF
	refTrx.RefID = reward[0].RefTrx
	refTrx.UsedReferralCode = rwdTrx.UsedPromoCode
	refTrx.RewardType = reward[0].Type
	refTrx.Type = ReferralType[reward[0].Reference]
	refTrx.RewardReferral = reward[0].Value
	refTrx.PhoneNumber = rwdTrx.RequestData.Phone

	return refTrx
}
