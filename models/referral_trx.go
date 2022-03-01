package models

import (
	"time"
)

var (
	// ReferralTrxTypeReferral to store referral_trx type referral
	ReferralTrxTypeReferral int64

	// ReferralTrxTypeReferrer to store referral_trx type referral
	ReferralTrxTypeReferrer int64 = 1

	// ReferralGoldback to store referral reward type goldback
	ReferralGoldback = "goldback"
)

// ReferralType to map referral type string to int64
var ReferralType = map[string]int64{
	RefTargetReferral: ReferralTrxTypeReferral,
	RefTargetReferrer: ReferralTrxTypeReferrer,
}

// ReferralTrx is represent a referral_transactions model
type ReferralTrx struct {
	ID               int64      `json:"id,omitempty"`
	CIF              string     `json:"cif,omitempty"`
	TotalGoldback    float64    `json:"totalGoldback,omitempty"`
	RefID            string     `json:"refId,omitempty"`
	UsedReferralCode string     `json:"usedReferralCode,omitempty"`
	RewardType       string     `json:"rewardType,omitempty"`
	PhoneNumber      string     `json:"phoneNumber,omitempty"`
	Type             int64      `json:"type,omitempty"`
	RewardReferral   int64      `json:"rewardReferral,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	CifReferrer      string     `json:"cifReferrer,omitempty"`
}

// Milestone is represent a referral_transactions model
type Milestone struct {
	Stages             int64   `json:"stage,omitempty"`
	LimitRewardCounter int64   `json:"limitRewardCounter,omitempty"`
	LimitReward        int64   `json:"limitReward,omitempty"`
	TotalRewardCounter int64   `json:"totalRewardCounter,omitempty"`
	TotalReward        int64   `json:"totalReward,omitempty"`
	Ranking            Ranking `json:"ranking,omitempty"`
}

// MilestonePayload is represent a referral_transactions model
type MilestonePayload struct {
	ReferralCode string `json:"referralCode,omitempty" validate:"required"`
}

// Ranking is represent a referral_transactions model
type Ranking struct {
	NoRanking      string `json:"noRanking,omitempty"`
	ReferralCode   string `json:"referralCode,omitempty"`
	TotalUsed      string `json:"totalUsed,omitempty"`
	Date           string `json:"date,omitempty"`
	IsReferralCode string `json:"isReferralCode,omitempty"`
}

// RankingPayload is represent a referral_transactions model
type RankingPayload struct {
	ReferralCode string `json:"referralCode,omitempty"`
	StartDate    string `json:"startDate,omitempty" validate:"required"`
	EndDate      string `json:"endDate,omitempty" validate:"required"`
}

type CoreTrxPayload struct {
	CIF              string  `json:"cif,omitempty"`
	RefID            string  `json:"refId,omitempty"`
	UsedReferralCode string  `json:"usedReferralCode,omitempty"`
	Type             int64   `json:"type,omitempty"`
	RewardReferral   int64   `json:"rewardReferral,omitempty"`
	RewardType       string  `json:"rewardType,omitempty"`
	PhoneNumber      string  `json:"phoneNumber,omitempty"`
	TrxAmount        float64 `json:"trxAmount,omitempty"`
	LoanAmount       float64 `json:"loanAmount,omitempty"`
	InterestAmount   float64 `json:"interestAmount,omitempty"`
	TrxDate          string  `json:"trxDate,omitempty"`
	ProductCode      int64   `json:"productCode,omitempty"`
}

type CoreTrxResponse struct {
	StatusCode int64  `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
}
