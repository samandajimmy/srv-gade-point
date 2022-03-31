package repository

import (
	"gade/srv-gade-point/api"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"github.com/labstack/echo"
)

type RestReferall struct {
	Xpoin api.IApiXpoin
}

func NewRestReferall(xpoin api.IApiXpoin) referrals.RestRefRepository {
	return &RestReferall{
		Xpoin: xpoin,
	}
}

func (rr *RestReferall) RRGetOslStatus(c echo.Context, pl models.ReqOslStatus) (bool, error) {
	resp, err := rr.Xpoin.XpoinPost(c, pl, "/xpoin/cgc/inactivecif")

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return false, models.ErrXpoinApi
	}

	if resp.ResponseCode == api.XpoinCodeOslInactive {
		return true, nil
	}

	return false, nil
}
