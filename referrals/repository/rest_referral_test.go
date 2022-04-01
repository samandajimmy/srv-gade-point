package repository_test

import (
	"gade/srv-gade-point/api"
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	refRepo "gade/srv-gade-point/referrals/repository"
	"gade/srv-gade-point/test"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	_ = BeforeSuite(func() {
		config.LoadEnv()
		config.LoadTestData()
	})

	_ = AfterSuite(func() {
	})
)

var _ = Describe("RestReferral", func() {
	var expectedResp api.XpoinResponse
	var mockCtrl *gomock.Controller
	var e test.DummyEcho
	var usedCif string
	var apiXpoin *mocks.MockIApiXpoin
	var obj models.ReqOslStatus
	var oslStatus bool
	var err error
	var endpoint string

	BeforeEach(func() {
		apiXpoin, expectedResp, endpoint = nil, api.XpoinResponse{}, ""
		usedCif = gofakeit.Regex("[1234567890]{10}")
		mockCtrl = gomock.NewController(GinkgoT())
		e = test.NewDummyEcho(http.MethodPost, "/", nil)
	})

	Describe("RRGetOslStatus", func() {
		BeforeEach(func() {
			endpoint = "/xpoin/cgc/inactivecif"
		})

		JustBeforeEach(func() {
			rr := refRepo.NewRestReferall(apiXpoin)
			oslStatus, err = rr.RRGetOslStatus(e.Context, obj)
		})

		Context("cif osl inactive", func() {
			BeforeEach(func() {
				expectedResp = api.XpoinResponse{
					ResponseCode: "00",
					ResponseDesc: "Approved",
				}
				obj = models.ReqOslStatus{
					Cif:         usedCif,
					ProductCode: "02",
				}

				apiXpoin = mocks.NewMockIApiXpoin(mockCtrl)
				apiXpoin.EXPECT().XpoinPost(e.Context, obj, endpoint).Return(expectedResp, nil)
			})

			It("expected to return success", func() {
				Expect([]interface{}{oslStatus, err}).To(Equal([]interface{}{true, nil}))
			})
		})

		Context("cif osl active", func() {
			BeforeEach(func() {
				expectedResp = api.XpoinResponse{
					ResponseCode: "404",
					ResponseDesc: "Rejected",
				}
				obj = models.ReqOslStatus{
					Cif:         usedCif,
					ProductCode: "02",
				}

				apiXpoin = mocks.NewMockIApiXpoin(mockCtrl)
				apiXpoin.EXPECT().XpoinPost(e.Context, obj, endpoint).Return(expectedResp, nil)
			})

			It("expected to return success", func() {
				Expect([]interface{}{oslStatus, err}).To(Equal([]interface{}{false, nil}))
			})
		})
	})
})
