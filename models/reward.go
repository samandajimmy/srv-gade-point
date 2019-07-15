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

var rewardType = map[int64]string{
	RewardTypePoint:    "point",
	RewardTypeDiscount: "discount",
	RewardTypeGoldback: "goldback",
	RewardTypeVoucher:  "voucher",
}

// Reward is represent a reward model
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
	CampaignID         *int64     `json:"campaign_id,omitempty"`
	Validators         *Validator `json:"validators,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
	Campaign           *Campaign  `json:"campaign,omitempty"`
	Quotas             *[]Quota   `json:"quotas,omitempty"`
	Tags               *[]Tag     `json:"tags,omitempty"`
	Vouchers           *[]Voucher `json:"vouchers,omitempty"`
}

// RewardsInquiry is represent a inquire reward response model
type RewardsInquiry struct {
	RefTrx  string            `json:"refTrx,omitempty"`
	Rewards *[]RewardResponse `json:"rewards,omitempty"`
}

// RewardResponse is represent a reward response model
type RewardResponse struct {
	Type           string  `json:"type,omitempty"`
	JournalAccount string  `json:"journalAccount,omitempty"`
	Value          float64 `json:"value,omitempty"`
	VoucherName    string  `json:"voucherName,omitempty"`
}

// GetRewardTypeText to get text of reward type
func (rwd Reward) GetRewardTypeText() string {
	return rewardType[*rwd.Type]
}
