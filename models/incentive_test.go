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

	Describe("ValidateMaxTransaction", func() {
		var output, amount float64
		incentive := models.Incentive{MaxTransaction: 1000}

		JustBeforeEach(func() {
			output = incentive.ValidateMaxTransaction(amount)
		})

		Context("when below the MaxTransaction", func() {
			BeforeEach(func() {
				amount = float64(200)
			})

			It("expect to return 200", func() {
				Expect(output).To(Equal(float64(200)))
			})
		})

		Context("when meet the MaxTransaction", func() {
			BeforeEach(func() {
				amount = float64(1000)
			})

			It("expect to return 1000", func() {
				Expect(output).To(Equal(float64(1000)))
			})
		})

		Context("when above the MaxTransaction", func() {
			BeforeEach(func() {
				amount = float64(2000)
			})

			It("expect to return 1000", func() {
				Expect(output).To(Equal(float64(1000)))
			})
		})
	})
})
