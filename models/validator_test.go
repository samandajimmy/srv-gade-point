package models_test

import (
	"errors"
	"gade/srv-gade-point/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	disc      float64 = 10
	maxValue  float64 = 15000
	minTrx    float64 = 100000
	trx       float64 = 1000000
	mltplr    float64 = 1
	val       float64 = 5000
	frml              = "(transactionAmount/multiplier)*value"
	voucherID int64   = 5
)

var validator = models.Validator{
	Channel:              "pds",
	Discount:             &disc,
	Formula:              frml,
	MaxValue:             &maxValue,
	MinTransactionAmount: &minTrx,
	Product:              "TE",
	TransactionType:      "OP",
	Multiplier:           &mltplr,
	Value:                &val,
	ValueVoucherID:       &voucherID,
}

var plValidator = models.PayloadValidator{
	CIF:               "1122334455",
	PromoCode:         "EMASMERDEKA",
	TransactionDate:   "2019-05-20",
	TransactionAmount: &trx,
	Validators: &models.Validator{
		Channel:         "pds",
		Product:         "TE",
		TransactionType: "OP",
	},
}

func TestValidate(t *testing.T) {
	// when its valid
	assert.Equal(t, nil, validator.Validate(&plValidator))

	// when its not valid
	plValidator.Validators.Channel = "cacing"
	err := errors.New("channel on this transaction is not valid to use the benefit")
	assert.Equal(t, err, validator.Validate(&plValidator))
}

func TestGetFormulaResult(t *testing.T) {
	// when its valid
	goldTrx := float64(5)
	plValidator.TransactionAmount = &goldTrx
	assert.Equal(t, mltplFunc(float64(25000), nil), mltplFunc(validator.GetFormulaResult(&plValidator)))

	// when formula is not available
	validator.Formula = ""
	assert.Equal(t, mltplFunc(float64(0), nil), mltplFunc(validator.GetFormulaResult(&plValidator)))
}

func TestCalculateDiscount(t *testing.T) {
	// when its valid
	hundredK := float64(100000)
	assert.Equal(t, mltplFunc(hundredK, nil), mltplFunc(validator.CalculateDiscount(plValidator.TransactionAmount)))

	// when discount is not available
	zero := float64(0)
	validator.Discount = &zero
	assert.Equal(t, mltplFunc(float64(0), nil), mltplFunc(validator.CalculateDiscount(plValidator.TransactionAmount)))
}

func TestGetVoucherResult(t *testing.T) {
	// when its valid
	assert.Equal(t, mltplFunc(int64(5), nil), mltplFunc(validator.GetVoucherResult()))

	// when value voucher ID is not available
	zero := int64(0)
	validator.ValueVoucherID = &zero
	assert.Equal(t, mltplFunc(int64(0), nil), mltplFunc(validator.GetVoucherResult()))
}

func TestCalculateMaximumValue(t *testing.T) {
	// when its valid
	assert.Equal(t, mltplFunc(float64(5000), nil), mltplFunc(validator.CalculateMaximumValue(&val)))

	// when maximum value is not available
	zero := float64(0)
	validator.MaxValue = &zero
	assert.Equal(t, mltplFunc(zero, nil), mltplFunc(validator.CalculateMaximumValue(&val)))
}

func TestGetRewardValue(t *testing.T) {
	// when its a direct discount reward
	zero := float64(0)
	validator.Formula = ""
	validator.Multiplier = &zero
	validator.MaxValue = &zero
	validator.Discount = &zero
	assert.Equal(t, mltplFunc(float64(5000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))

	// when its a discount percentage reward
	validator.Discount = &disc
	validator.Value = &zero
	assert.Equal(t, mltplFunc(float64(100000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))

	// when its a discount percentage reward with nil value
	validator.Discount = &disc
	validator.Value = nil
	assert.Equal(t, mltplFunc(float64(100000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))

	// when its a discount percentage with max value reward
	validator.MaxValue = &maxValue
	assert.Equal(t, mltplFunc(float64(15000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))

	// when its a direct discount and discount percentage reward
	validator.MaxValue = &zero
	validator.Discount = &disc
	validator.Value = &val
	assert.Equal(t, mltplFunc(float64(105000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))

	// when its a multiplied discount reward
	five := float64(5)
	one := float64(1)
	validator.MaxValue = &zero
	validator.Discount = &zero
	validator.Formula = frml
	validator.Multiplier = &mltplr
	validator.MinTransactionAmount = &one
	plValidator.TransactionAmount = &five
	assert.Equal(t, mltplFunc(float64(25000), nil), mltplFunc(validator.GetRewardValue(&plValidator)))
}

func mltplFunc(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}
