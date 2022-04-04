package models_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
)

var _ = Describe("Validator", func() {

	Describe("GetRewardValue", func() {
		Context("minimum direct discount reward", func() {
			validator := models.Validator{
				Channel: "9997", Product: "62", TransactionType: "OP",
				Value: helper.CreateFloat64(50000)}

			plValidator := models.PayloadValidator{
				TransactionAmount: helper.CreateFloat64(100000),
				Validators: &models.Validator{
					Channel: "9997", Product: "62", TransactionType: "OP"}}

			value, err := validator.GetRewardValue(&plValidator)

			It("happy test", func() {
				Expect(value).To(Equal(float64(50000)))
				Expect(err).To(BeNil())
			})
		})

		Context("minimum percentage discount reward", func() {
			validator := models.Validator{
				Channel: "9997", Product: "62", TransactionType: "OP",
				Discount: helper.CreateFloat64(20)}

			plValidator := models.PayloadValidator{
				TransactionAmount: helper.CreateFloat64(400000),
				Validators: &models.Validator{
					Channel: "9997", Product: "62", TransactionType: "OP"}}

			value, err := validator.GetRewardValue(&plValidator)

			It("happy test", func() {
				Expect(value).To(Equal(float64(80000)))
				Expect(err).To(BeNil())
			})
		})

		Context("minimum goldback reward", func() {
			validator := models.Validator{
				Channel: "9997", Product: "62", TransactionType: "OP",
				Multiplier: helper.CreateFloat64(10000),
				Value:      helper.CreateFloat64(2),
				Formula:    "(transactionAmount/multiplier)*value"}

			plValidator := models.PayloadValidator{
				TransactionAmount: helper.CreateFloat64(50000),
				Validators: &models.Validator{
					Channel: "9997", Product: "62", TransactionType: "OP"}}

			value, err := validator.GetRewardValue(&plValidator)

			It("happy test", func() {
				Expect(value).To(Equal(float64(10)))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Validate", func() {
		validator := models.Validator{
			Channel: "9997", Product: "62", TransactionType: "OP",
			MinTransactionAmount: helper.CreateFloat64(100000),
			Value:                helper.CreateFloat64(50000)}

		plValidator := models.PayloadValidator{
			TransactionAmount: helper.CreateFloat64(100000),
			Validators: &models.Validator{
				Channel: "9997", Product: "62", TransactionType: "OP"}}

		err := validator.Validate(&plValidator)

		It("happy test", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("GetFormulaResult", func() {
		var validator models.Validator
		var plValidator models.PayloadValidator
		var value float64
		var err error

		BeforeEach(func() {
			plValidator = models.PayloadValidator{TransactionAmount: helper.CreateFloat64(500000)}
			validator = models.Validator{
				Multiplier: helper.CreateFloat64(10000),
				Value:      helper.CreateFloat64(100),
				Formula:    "(transactionAmount/multiplier)*value",
			}
			value = 0
			err = nil
		})

		JustBeforeEach(func() {
			value, err = validator.GetFormulaResult(&plValidator)
		})

		Context("when formula expresion is an empty string", func() {

			BeforeEach(func() {
				validator.Formula = ""
			})

			It("expect return value is 0", func() {
				Expect(value).To(Equal(float64(0)))

			})

			It("expect return error is nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when formula expresion can't be converted", func() {

			BeforeEach(func() {
				validator.Formula = "Test 1234 ___ !!!"
			})

			It("expect return value is 0", func() {
				Expect(value).To(Equal(float64(0)))

			})

			It("expect formula conversion to return error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when formula didn't match the payload", func() {

			BeforeEach(func() {
				validator.Formula = "(transactionAmount/multiplier)*value/divider"
			})

			It("expect return value is 0", func() {
				Expect(value).To(Equal(float64(0)))

			})

			It("expect formula conversion to return setup error", func() {
				Expect(err).To(MatchError(models.ErrFormulaSetup))
			})
		})

		Context("when GetFormulaResult return result", func() {

			BeforeEach(func() {
				validator.Formula = "(transactionAmount/multiplier)*value"
			})

			It("expect return value is 5000", func() {
				Expect(value).To(Equal(float64(5000)))
			})

			It("expect GetFormulaResult return error is nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when parameters are from payload validator", func() {

			BeforeEach(func() {
				plValidator.TransactionAmount = helper.CreateFloat64(750000)
				plValidator.LoanAmount = helper.CreateFloat64(250000)
				validator.Formula = "(transactionAmount/loanAmount)"
			})

			It("expect return value is 3", func() {
				Expect(value).To(Equal(float64(3)))
			})

			It("expect GetFormulaResult return error is nil", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetVoucherResult", func() {
		validator := models.Validator{
			ValueVoucherID: helper.CreateInt64(2),
		}

		value, err := validator.GetVoucherResult()

		It("happy test", func() {
			Expect(value).To(Equal(int64(2)))
			Expect(err).To(BeNil())
		})
	})

	Describe("CalculateDiscount", func() {
		validator := models.Validator{
			Discount: helper.CreateFloat64(10),
		}

		value, err := validator.CalculateDiscount(helper.CreateFloat64(50000))

		It("happy test", func() {
			Expect(value).To(Equal(float64(5000)))
			Expect(err).To(BeNil())
		})
	})

	Describe("CalculateMaximumValue", func() {
		validator := models.Validator{
			MaxValue: helper.CreateFloat64(40000),
		}

		value, err := validator.CalculateMaximumValue(helper.CreateFloat64(50000))

		It("happy test", func() {
			Expect(value).To(Equal(float64(40000)))
			Expect(err).To(BeNil())
		})
	})

})
