package models

import "github.com/dgrijalva/jwt-go"

// Token to store JWT token data
type Token struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// AccountToken to store account token data
type AccountToken struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}
