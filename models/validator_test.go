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
		validator := models.Validator{
			Multiplier: helper.CreateFloat64(10000),
			Value:      helper.CreateFloat64(100),
			Formula:    "(transactionAmount/multiplier)*value",
		}

		plValidator := models.PayloadValidator{TransactionAmount: helper.CreateFloat64(500000)}

		value, err := validator.GetFormulaResult(&plValidator)

		It("happy test", func() {
			Expect(value).To(Equal(float64(5000)))
			Expect(err).To(BeNil())
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
