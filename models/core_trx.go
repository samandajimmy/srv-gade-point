package models

type CoreTrxPayload struct {
	CIF            string  `json:"cif"`
	PhoneNumber    string  `json:"phoneNumber,omitempty"`
	TrxID          string  `json:"transactionId"`
	TrxAmount      float64 `json:"transactionAmount"`
	LoanAmount     float64 `json:"loanAmount"`
	InterestAmount float64 `json:"interestAmount"`
	MarketingCode  string  `json:"marketingCode"`
	TrxDate        string  `json:"transactionDate"`
	ProductCode    int64   `json:"productCode"`
}

type CoreTrxResponse struct {
	StatusCode int64  `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
}
