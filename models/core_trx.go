package models

type CoreTrxPayload struct {
	CIF            string  `json:"cif"`
	PhoneNumber    string  `json:"phoneNumber,omitempty"`
	TrxID          string  `json:"trxId"`
	TrxAmount      float64 `json:"trxAmount"`
	LoanAmount     float64 `json:"loanAmount"`
	InterestAmount float64 `json:"interestAmount"`
	MarketingCode  string  `json:"marketingCode"`
	TrxDate        string  `json:"trxDate"`
	ProductCode    int64   `json:"productCode"`
}

type CoreTrxResponse struct {
	StatusCode int64  `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
}
