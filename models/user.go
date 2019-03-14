package models

import "time"

// Validators is represent a validators inside user
type ValidatorsUser struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// User is represent a users model
type User struct {
	ID        int64      `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	Status    *int8      `json:"status,omitempty"`
	Role      *int8      `json:"role,omitempty"`
	Email     string     `json:"email,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}
