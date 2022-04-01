package models

var (
	CoreTrxSuccess       = "success"
	CoreTrxInqFailed     = "failed_inquiry"
	CoreTrxPaymentFailed = "failed_payment"
)

type CoreTrxPayload struct {
	CIF            string  `json:"cif"`
	PhoneNumber    string  `json:"phoneNumber,omitempty"`
	CustomerName   string  `json:"customerName,omitempty"`
	TrxID          string  `json:"trxId"`
	TrxAmount      float64 `json:"trxAmount"`
	LoanAmount     float64 `json:"loanAmount"`
	InterestAmount float64 `json:"interestAmount"`
	MarketingCode  string  `json:"marketingCode"`
	TrxDate        string  `json:"trxDate"`
	ProductCode    string  `json:"productCode"`
	TrxType        string  `json:"trxType"`
	Channel        string  `json:"channel,omitempty"`
	InqStatus      string  `json:"inqStatus"`
	RefTrx         string  `json:"refTrx"`
}

type CoreTrxResponse struct {
	StatusCode int64  `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
}
