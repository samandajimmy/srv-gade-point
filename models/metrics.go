package models

import "time"

// Metrics is represent a metrics model
type Metrics struct {
	ID               int64      `json:"id,omitempty"`
	Job              string     `json:"job,omitempty"`
	Counter          *int64     `json:"counter,omitempty"`
	Status           string     `json:"status,omitempty"`
	CreationTime     *time.Time `json:"creation_time,omitempty"`
	ModificationTime *time.Time `json:"modification_time,omitempty"`
}
