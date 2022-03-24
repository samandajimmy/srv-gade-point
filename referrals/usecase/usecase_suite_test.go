package usecase_test

import (
	"gade/srv-gade-point/config"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	config.LoadEnv()
	config.LoadTestData()
})

func TestUsecase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Referral Usecase Suite")
}
