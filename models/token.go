package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token to store JWT token data
type Token struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// AccountToken to store account token data
type AccountToken struct {
	ID        int64      `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	Token     string     `json:"token,omitempty"`
	ExpireAt  *time.Time `json:"expireAt,omitempty"`
	Status    *int8      `json:"status,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
