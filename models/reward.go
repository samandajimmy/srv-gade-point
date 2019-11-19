package models

import (
	"encoding/json"
	"math"
	"time"
)

var (
	// RewardTypePoint to store reward type point
	RewardTypePoint int64
	// RewardTypeDirectDiscount to store reward type direct discount
	RewardTypeDirectDiscount int64 = 1
	// RewardTypePercentageDiscount to store reward type percentage discount
	RewardTypePercentageDiscount int64 = 2
	// RewardTypeGoldback to store reward type goldback
	RewardTypeGoldback int64 = 3
	// RewardTypeVoucher to store reward type voucher
	RewardTypeVoucher int64 = 4

	// IsPromoCodeFalse to store is promo code false
	IsPromoCodeFalse int64
	// IsPromoCodeTrue to store is promo code true
	IsPromoCodeTrue int64 = 1
)

var rewardType = map[int64]string{
	RewardTypePoint:              "point",
	RewardTypeDirectDiscount:     "discount",
	RewardTypePercentageDiscount: "discount",
	RewardTypeGoldback:           "goldback",
	RewardTypeVoucher:            "voucher",
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
	RefID              string     `json:"-"`
	RootRefID          string     `json:"-"`
}

// RewardsInquiry is represent a inquire reward response model
type RewardsInquiry struct {
	RefTrx  string            `json:"refTrx,omitempty"`
	Rewards *[]RewardResponse `json:"rewards,omitempty"`
}

// RewardResponse is represent a reward response model
type RewardResponse struct {
	RewardID       int64   `json:"-"`
	CIF            string  `json:"cif,omitempty"`
	Product        string  `json:"product,omitempty"`
	RefTrx         string  `json:"refTrx,omitempty"`
	Type           string  `json:"type,omitempty"`
	JournalAccount string  `json:"journalAccount,omitempty"`
	Value          float64 `json:"value,omitempty"`
	VoucherName    string  `json:"voucherName,omitempty"`
}

// RewardPayment is represent a reward payment model
type RewardPayment struct {
	CIF     string `json:"cif,omitempty"`
	RefTrx  string `json:"refTrx,omitempty"`
	RefCore string `json:"refCore,omitempty"`
}

// RewardsPayload is represent a inquire reward and reward transactiomn response model
type RewardsPayload struct {
	RewardID             string `json:"rewardID,omitempty"`
	Status               string `json:"status,omitempty"`
	Page                 int    `json:"page,omitempty"`
	Limit                int    `json:"limit,omitempty"`
	StartTransactionDate string `json:"startTransactionDate,omitempty"`
	EndTransactionDate   string `json:"endTransactionDate,omitempty"`
	StartSuccededDate    string `json:"startSuccededDate,omitempty"`
	EndSuccededDate      string `json:"endSuccededDate,omitempty"`
}

// GetRewardTypeText to get text of reward type
func (rwd Reward) GetRewardTypeText() string {
	return rewardType[*rwd.Type]
}

// Populate to generate a reward response
func (rwdResp *RewardResponse) Populate(reward Reward, rwdValue float64, pv PayloadValidator) {
	rwdResp.Type = reward.GetRewardTypeText()
	rwdResp.Value = roundDown(rwdValue, 0)
	rwdResp.JournalAccount = reward.JournalAccount
	rwdResp.RefTrx = reward.RefID
	rwdResp.RewardID = reward.ID
	rwdResp.CIF = pv.CIF
	rwdResp.Product = pv.Validators.Product

	// convert pv struct to map
	pvMap := map[string]interface{}{}
	tempJSON, _ := json.Marshal(pv)
	json.Unmarshal(tempJSON, &pvMap)

	// validate referral reward
	if referrer := reward.Validators.ReferralTarget(pv); referrer != "" {
		rwdResp.CIF = pvMap[referrer].(string)
	}
}

func roundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
