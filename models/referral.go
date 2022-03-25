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

type SumIncentive struct {
	PerDay        float64 `json:"perDay"`
	PerMonth      float64 `json:"perMonth"`
	Total         float64 `json:"total"`
	ValidPerDay   bool    `json:"validPerDay"`
	ValidPerMonth bool    `json:"validPerMonth"`
	IsValid       bool    `json:"isValid"`
	Reward        Reward  `json:"-"`
}

type PrefixResponse struct {
	Prefix string `json:"prefix,omitempty"`
}

// RespFriends list member used code referral
type RespFriends struct {
	CustomerName string `json:"customerName,omitempty"`
}

type PayloadFriends struct {
	CIF   string `json:"cif,omitempty"`
	Limit int8   `json:"limit"`
	Page  int8   `json:"page"`
}
