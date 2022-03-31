package models_test

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Incentive", func() {
	Describe("ValidateMaxIncentive", func() {
		var incentive models.Incentive
		var obj models.ObjIncentive

		JustBeforeEach(func() {
			incentive.Validate(&obj)
		})

		BeforeEach(func() {
			incentive = models.Incentive{}
			obj = models.ObjIncentive{}
		})

		Context(`incentive per day greater than max per day
			and incentive per month greater than max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = helper.CreateFloat64(500)
				incentive.MaxPerMonth = helper.CreateFloat64(500)
				obj.PerDay = float64(1000)
				obj.PerMonth = float64(1000)
			})

			It("expect to have incentive per day equal to max per day", func() {
				Expect(obj.PerDay).To(Equal(*incentive.MaxPerDay))
			})

			It("expect incentive validPerDay to be false", func() {
				Expect(obj.ValidPerDay).To(BeFalse())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(obj.PerDay).To(Equal(*incentive.MaxPerDay))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(obj.ValidPerMonth).To(BeFalse())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(obj.IsValid).To(BeFalse())
			})
		})

		Context(`incentive per day less than max per day
			and incentive per month less than max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = helper.CreateFloat64(1000)
				incentive.MaxPerMonth = helper.CreateFloat64(1000)
				obj.PerDay = float64(100)
				obj.PerMonth = float64(100)
			})

			It("expect to have incentive per day equal to 100", func() {
				Expect(obj.PerDay).To(Equal(float64(100)))
			})

			It("expect incentive validPerDay to be true", func() {
				Expect(obj.ValidPerDay).To(BeTrue())
			})

			It("expect to have incentive per month equal to 100", func() {
				Expect(obj.PerDay).To(Equal(float64(100)))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(obj.ValidPerMonth).To(BeTrue())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(obj.IsValid).To(BeTrue())
			})
		})

		Context(`incentive per day equal to max per day
			and incentive per month equal to max per month`, func() {
			BeforeEach(func() {
				incentive.MaxPerDay = helper.CreateFloat64(1000)
				incentive.MaxPerMonth = helper.CreateFloat64(1000)
				obj.PerDay = float64(1000)
				obj.PerMonth = float64(1000)
			})

			It("expect to have incentive per day equal to 1000", func() {
				Expect(obj.PerDay).To(Equal(float64(1000)))
			})

			It("expect incentive validPerDay to be true", func() {
				Expect(obj.ValidPerDay).To(BeTrue())
			})

			It("expect to have incentive per month equal to 1000", func() {
				Expect(obj.PerDay).To(Equal(float64(1000)))
			})

			It("expect incentive validPerMonth to be false", func() {
				Expect(obj.ValidPerMonth).To(BeTrue())
			})

			It("expect to have incentive per month equal to max per month", func() {
				Expect(obj.IsValid).To(BeTrue())
			})
		})

		Context("ValidateMaxTransaction", func() {
			incentive := models.Incentive{MaxPerTransaction: helper.CreateFloat64(1000)}

			JustBeforeEach(func() {
				incentive.Validate(&obj)
			})

			Context("when below the MaxTransaction", func() {
				BeforeEach(func() {
					obj.PerTransaction = float64(200)
				})

				It("expect to return 200", func() {
					Expect(obj.ValidTransaction).To(BeTrue())
				})
			})

			Context("when meet the MaxTransaction", func() {
				BeforeEach(func() {
					obj.PerTransaction = float64(1000)
				})

				It("expect to return 1000", func() {
					Expect(obj.ValidTransaction).To(BeTrue())
				})
			})

			Context("when above the MaxTransaction", func() {
				BeforeEach(func() {
					obj.PerTransaction = float64(2000)
				})

				It("expect to return 1000", func() {
					Expect(obj.ValidTransaction).To(BeFalse())
				})
			})
		})
	})
})
