package models

type PayloadVoucherValidator struct {
	StartDate       string `json:"startDate, omitempty"`
	EndDate         string `json:"endDate, omitempty"`
	ProductCode     string `json:"productCode", omitempty`
	TransactionType string `json:"transactionType", omitempty`
	Page            string `json:"page, omitempty" validate:"required"`
	Limit           string `json:"limit, omitempty validate:"required"`
}
