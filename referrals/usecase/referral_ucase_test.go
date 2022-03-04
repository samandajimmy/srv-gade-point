package usecase_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/referrals/repository"
	"gade/srv-gade-point/referrals/usecase"
)

var _ = Describe("ReferralUcase", func() {
	It("happy test", func() {
		config.LoadEnv()
		sqldb, dbBun := database.GetDBConn()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := repository.NewPsqlReferralRepository(sqldb, dbBun)
		rUS := usecase.NewReferralUseCase(repo)

		_, err := rUS.UMaxIncentiveValidation(c, "PSESOME")

		Expect(err).To(BeNil())
	})
})
