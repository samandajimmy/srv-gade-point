package models

import (
	"time"
)

var (
	// RewardTypePoint to store reward type point
	RewardTypePoint int64
	// RewardTypeDiscount to store reward type discount
	RewardTypeDiscount int64 = 1
	// RewardTypeGoldback to store reward type goldback
	RewardTypeGoldback int64 = 2
	// RewardTypeVoucher to store reward type voucher
	RewardTypeVoucher int64 = 3

	// IsPromoCodeFalse to store is promo code false
	IsPromoCodeFalse int64
	// IsPromoCodeTrue to store is promo code true
	IsPromoCodeTrue int64 = 1
)

// Reward is represent a point model
type Reward struct {
	ID                 int64      `json:"id,omitempty"`
	Name               string     `json:"name,omitempty"`
	Description        string     `json:"description,omitempty"`
	TermsAndConditions string     `json:"termsAndConditions,omitempty"`
	HowToUse           string     `json:"howToUse,omitempty"`
	PromoCode          string     `json:"promoCode,omitempty"`
	CustomPeriod       string     `json:"customPeriod,omitempty"`
	JournalAccount     string     `json:"journalAccount,omitempty"`
	IsPromoCode        *int64     `json:"isPromoCode,omitempty"`
	Type               *int64     `json:"type,omitempty"`
	Validators         *Validator `json:"validators,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
	Campaign           *Campaign  `json:"campaign,omitempty"`
	Quotas             *[]Quota   `json:"quotas,omitempty"`
	Vouchers           *[]Voucher `json:"vouchers,omitempty"`
}
