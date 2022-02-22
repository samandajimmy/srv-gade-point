package repository

import (
	"database/sql"
	"gade/srv-gade-point/referrals"
)

type psqlReferralsRepository struct {
	Conn *sql.DB
}

// NewPsqlReferralRepository will create an object that represent the referrals.Repository interface
func NewPsqlReferralRepository(Conn *sql.DB) referrals.Repository {
	return &psqlReferralsRepository{Conn}
}
