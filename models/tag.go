package models

import (
	"time"
)

// Tag is represent a quota model
type Tag struct {
	ID        int64      `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}
