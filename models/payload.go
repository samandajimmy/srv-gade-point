package models

// PayloadList a struct to store all payload for a list response
type PayloadList struct {
	Status          string `json:"status,omitempty"`
	StartDate       string `json:"startDate,omitempty" validate:"dateString"`
	EndDate         string `json:"endDate,omitempty"`
	ProductCode     string `json:"productCode,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	Page            int64  `json:"page,omitempty" validate:"min=1"`
	Limit           int64  `json:"limit,omitempty" validate:"min=1"`
}

// RewardPromotionLists a struct to store payload RewardPromotionLists for a list response
type RewardPromotionLists struct {
	Product         string `json:"product,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	Channel         string `json:"channel,omitempty"`
}

type RespReferral struct {
	CIF          string `json:"cif,omitempty"`
	ReferralCode string `json:"referralCode,omitempty"`
}

type RespReferralDetail struct {
	ReferralCode string       `json:"referralCode"`
	Incentive    SumIncentive `json:"incentive"`
}

type RequestHistoryIncentive struct {
	RefCif string `json:"refCif" validate:"required"`
	Limit  int64  `json:"limit" validate:"required"`
	Page   int64  `json:"page"`
}

type ResponseHistoryIncentive struct {
	TotalData            int64                           `json:"totalData"`
	HistoryIncentiveData *[]ResponseHistoryIncentiveData `json:"historyIncentiveData"`
}

type ResponseHistoryIncentiveData struct {
	TransactionType string  `json:"transactionType"`
	ProductCode     string  `json:"productCode"`
	CustomerName    string  `json:"customerName"`
	RewardReferral  float64 `json:"rewardReferral"`
	CreatedAt       float64 `json:"createdAt"`
}
