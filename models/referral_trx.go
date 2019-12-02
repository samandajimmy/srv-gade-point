package models

import "time"

var (
	// ReferralTrxTypeReferral to store referral_trx type referral
	ReferralTrxTypeReferral int64

	// ReferralTrxTypeReferrer to store referral_trx type referral
	ReferralTrxTypeReferrer int64 = 1
)

// ReferralType to map referral type string to int64
var ReferralType = map[string]int64{
	CampaignCodeReferral: ReferralTrxTypeReferral,
	CampaignCodeReferrer: ReferralTrxTypeReferrer,
}

// ReferralTrx is represent a referral_transactions model
type ReferralTrx struct {
	ID               int64      `json:"id,omitempty"`
	CIF              string     `json:"cif,omitempty"`
	RefID            string     `json:"refId,omitempty"`
	UsedReferralCode string     `json:"usedReferralCode,omitempty"`
	RewardType       string     `json:"rewardType,omitempty"`
	Type             int64      `json:"type,omitempty"`
	RewardReferral   int64      `json:"rewardReferral,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}
