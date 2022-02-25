package models

import "time"

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
	CIF    string `json:"cif,omitempty" validate:"required"`
	Prefix string `json:"prefix,omitempty" validate:"required,max=5"`
}

type ReferralCodeUser struct {
	ReferralCode string `json:"referralCode,omitempty"`
}

type RequestReferralCodeUser struct {
	CIF string `json:"cif,omitempty" validate:"required"`
}
