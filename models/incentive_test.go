package models_test

import (
	"gade/srv-gade-point/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Incentive", func() {
	Describe("ValidateMaxIncentive", func() {
		var incentive models.Incentive
		var sum models.SumIncentive

		JustBeforeEach(func() {
			incentive.ValidateMaxIncentive(&sum)
		})

		BeforeEach(func() {
			incentive = models.Incentive{}
			sum = models.SumIncentive{}
		})

		Context(`incentive per day greater than max per day
			and incentive per month greater than max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = float64(500)
				incentive.MaxPerMonth = float64(500)
				sum.PerDay = float64(1000)
				sum.PerMonth = float64(1000)
			})

			It("expect to have incentive per day equal to max per day", func() {
				Expect(sum.PerDay).To(Equal(incentive.MaxPerDay))
			})

			It("expect incentive validPerDay to be false", func() {
				Expect(sum.ValidPerDay).To(BeFalse())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(sum.PerDay).To(Equal(incentive.MaxPerDay))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(sum.ValidPerMonth).To(BeFalse())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(sum.IsValid).To(BeFalse())
			})
		})

		Context(`incentive per day less than max per day
			and incentive per month less than max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = float64(1000)
				incentive.MaxPerMonth = float64(1000)
				sum.PerDay = float64(100)
				sum.PerMonth = float64(100)
			})

			It("expect to have incentive per day equal to 100", func() {
				Expect(sum.PerDay).To(Equal(float64(100)))
			})

			It("expect incentive validPerDay to be true", func() {
				Expect(sum.ValidPerDay).To(BeTrue())
			})

			It("expect to have incentive per month equal to 100", func() {
				Expect(sum.PerDay).To(Equal(float64(100)))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(sum.ValidPerMonth).To(BeTrue())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(sum.IsValid).To(BeTrue())
			})
		})

		Context(`incentive per day equal to max per day
			and incentive per month equal to max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = float64(1000)
				incentive.MaxPerMonth = float64(1000)
				sum.PerDay = float64(1000)
				sum.PerMonth = float64(1000)
			})

			It("expect to have incentive per day equal to 1000", func() {
				Expect(sum.PerDay).To(Equal(float64(1000)))
			})

			It("expect incentive validPerDay to be true", func() {
				Expect(sum.ValidPerDay).To(BeTrue())
			})

			It("expect to have incentive per month equal to 1000", func() {
				Expect(sum.PerDay).To(Equal(float64(1000)))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(sum.ValidPerMonth).To(BeTrue())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(sum.IsValid).To(BeTrue())
			})
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
