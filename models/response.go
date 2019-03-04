package models

// Response struct is represent a data for output payload
type Response struct {
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	TotalCount string      `json:"totalCount,omitempty"`
}
