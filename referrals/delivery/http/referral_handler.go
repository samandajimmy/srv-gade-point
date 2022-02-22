package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
)

var response models.Response

// ReferralHandler represent the httphandler for referrals
type ReferralHandler struct {
	ReferralUseCase referrals.UseCase
}

// NewReferralsHandler represent to register referrals endpoint
func NewReferralsHandler(echoGroup models.EchoGroup, rus referrals.UseCase) {

}
