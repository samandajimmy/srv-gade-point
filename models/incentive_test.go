package models_test

import (
	"gade/srv-gade-point/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Incentive", func() {
	Describe("ValidateMaxIncentive", func() {
		It("happy test", func() {
			incentive := models.Incentive{
				MaxPerDay:   10,
				MaxPerMonth: 10,
			}
			sum := models.SumIncentive{
				PerDay:   100,
				PerMonth: 100,
			}

			incentive.ValidateMaxIncentive(&sum)

			Expect(sum).To(Equal(models.SumIncentive{
				PerDay:        10,
				PerMonth:      10,
				Total:         0,
				ValidPerDay:   false,
				ValidPerMonth: false,
			}))
		})
	})
})
