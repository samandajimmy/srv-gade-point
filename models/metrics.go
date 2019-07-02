package models

// Metrics is represent a metrics model
type Metrics struct {
	ID      int64  `json:"id,omitempty"`
	Module  string `json:"module,omitempty"`
	Counter *int64 `json:"status,omitempty"`
}
