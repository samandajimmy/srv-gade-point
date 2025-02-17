package models

import (
	"time"
)

// ReferralCodes is represent a referral_transactions model
type ReferralCodes struct {
	ID           int64     `json:"id,omitempty"`
	CIF          string    `json:"cif,omitempty" validate:"required"`
	ReferralCode string    `json:"referralCode,omitempty"`
	CampaignId   int64     `json:"campaignId,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

type RequestCreateReferral struct {
	CIF    string `json:"cif,omitempty" validate:"required,max=10"`
	Prefix string `json:"prefix,omitempty" validate:"required,max=5"`
}

type RequestReferralCodeUser struct {
	CIF string `json:"cif" validate:"required"`
}

type RefCampaignReward struct {
	CampaignId     int64     `json:"campaignId"`
	RewardId       int64     `json:"rewardId"`
	ReferralCodeId int64     `json:"referralCodeId"`
	EndDate        string    `json:"endDate"`
	StartDate      string    `json:"startDate"`
	Validators     Validator `json:"validators"`
	ReferralCode   string    `json:"referralCode"`
}

type PrefixResponse struct {
	Prefix string `json:"prefix,omitempty"`
}

type RespTotalFriends struct {
	TotalFriends int `json:"totalFriends"`
}

type Friends struct {
	CustomerName string `json:"customerName"`
}

type PayloadFriends struct {
	CIF   string `json:"cif,omitempty"`
	Limit int8   `json:"limit"`
	Page  int8   `json:"page"`
}
