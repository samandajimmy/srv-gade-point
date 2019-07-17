package models

import (
	"time"
)

var (
	// IsPerUserFalse to store is per user false
	IsPerUserFalse int64
	// IsPerUserTrue to store is per user true
	IsPerUserTrue int64 = 1

	//IsLimitAmount to store limit amount
	IsLimitAmount int64
)

// Quota is represent a quota model
type Quota struct {
	ID           int64      `json:"id,omitempty"`
	NumberOfDays *int64     `json:"numberOfDays,omitempty"`
	Amount       *int64     `json:"amount,omitempty"`
	IsPerUser    *int64     `json:"isPerUser,omitempty"`
	Available    *int64     `json:"available,omitempty"`
	LastCheck    *time.Time `json:"lastCheck,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	Reward       *Reward    `json:"reward,omitempty"`
}
