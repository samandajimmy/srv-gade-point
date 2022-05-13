package models

import (
	"fmt"
	"gade/srv-gade-point/logger"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	govaluate "gopkg.in/Knetic/govaluate.v2"
)

// Validator to store all validator data
type Validator struct {
	CampaignCode         string     `json:"campaignCode,omitempty"`
	Channel              string     `json:"channel,omitempty" validate:"required"`
	Discount             *float64   `json:"discount,omitempty"`
	Formula              string     `json:"formula,omitempty"`
	MaxLoanAmount        *float64   `json:"maxLoanAmount,omitempty"`
	MaxValue             *float64   `json:"maxValue,omitempty"`
	MinLoanAmount        *float64   `json:"minLoanAmount,omitempty"`
	MinTransactionAmount *float64   `json:"minTransactionAmount,omitempty"`
	Multiplier           *float64   `json:"multiplier,omitempty"`
	Deviden              *float64   `json:"deviden,omitempty"`
	Product              string     `json:"product,omitempty" validate:"required"`
	Source               string     `json:"source,omitempty"` // device name that user used
	TransactionType      string     `json:"transactionType,omitempty" validate:"required"`
	Unit                 string     `json:"unit,omitempty"`
	Target               string     `json:"target,omitempty"`
	IsPrivate            int        `json:"isPrivate"`
	IsDecimal            bool       `json:"isDecimal"`
	Value                *float64   `json:"value,omitempty"`
	ValueVoucherID       *int64     `json:"valueVoucherId,omitempty"`
	Incentive            *Incentive `json:"incentive,omitempty"`
}

// PayloadValidator to store a payload to validate a request
type PayloadValidator struct {
	BranchCode        string     `json:"branchCode,omitempty"`
	IsMulti           bool       `json:"isMulti,omitempty"`
	CampaignID        string     `json:"campaignId,omitempty"`
	CIF               string     `json:"cif,omitempty"`
	Referrer          string     `json:"referrer,omitempty" validate:"isRequiredWith=Validators.CampaignCode"`
	CustomerName      string     `json:"customerName,omitempty"`
	LoanAmount        *float64   `json:"loanAmount,omitempty"`
	Phone             string     `json:"phone,omitempty"`
	PromoCode         string     `json:"promoCode,omitempty" validate:"required"`
	RedeemedDate      string     `json:"redeemedDate,omitempty"`
	RefTrx            string     `json:"refTrx,omitempty"`
	TransactionDate   string     `json:"transactionDate,omitempty" validate:"required,dateString=2006-01-02 15:04:05.000"`
	TransactionAmount *float64   `json:"transactionAmount,omitempty" validate:"required"`
	InterestAmount    *float64   `json:"interestAmount,omitempty"`
	VoucherID         string     `json:"voucherId,omitempty"`
	Validators        *Validator `json:"validators,omitempty"`
	CodeType          string     `json:"codeType,omitempty"`
}

var skippedValidator = []string{"multiplier", "value", "formula", "maxValue", "unit", "isPrivate",
	"isDecimal", "incentive"}
var compareEqual = []string{"channel", "product", "transactionType", "source", "campaignCode"}
var evalFunc = []string{"floor"}
var tightenValidator = map[string]string{
	"minTransactionAmount": "transactionAmount",
	"minLoanAmount":        "loanAmount",
	"maxLoanAmount":        "loanAmount",
}
var customErrMsg = "%s untuk transaksi ini tidak valid untuk mendapatkan benefit"

// GetRewardValue is to get value of the reward
func (v *Validator) GetRewardValue(payloadValidator *PayloadValidator) (float64, error) {
	var result float64
	// validate eligibility
	err := v.Validate(payloadValidator)

	if err != nil {
		return 0, err
	}

	if v.Value != nil {
		result = *v.Value
	}

	// calculate formula if any
	formulaResult, err := v.GetFormulaResult(payloadValidator)

	if err != nil {
		return 0, err
	}

	if formulaResult != 0 {
		result = formulaResult
	}

	// calculate discount if any
	discResult, err := v.CalculateDiscount(payloadValidator.TransactionAmount)

	if err != nil {
		return 0, err
	}

	result = result + discResult

	// calculate maximum value if any
	MaxResult, err := v.CalculateMaximumValue(&result)

	if err != nil {
		return 0, err
	}

	if MaxResult != 0 {
		result = MaxResult
	}

	return result, nil
}

// Validate to validate client input with admin input
func (v *Validator) Validate(payloadValidator *PayloadValidator) error {
	ov := ObjectValidator{
		SkippedValidator: skippedValidator,
		CompareEqual:     compareEqual,
		TightenValidator: tightenValidator,
	}

	err := ov.autoValidating(v, payloadValidator, payloadValidator.Validators)

	if err != nil {
		logger.Make(nil).Debug(err)

		return err
	}

	return nil
}

// GetFormulaResult to proccess the formula then get the result
func (v *Validator) GetFormulaResult(payloadValidator *PayloadValidator) (float64, error) {

	// check formula availability
	if v.Formula == "" {
		return float64(0), nil
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(v.Formula, goValuateFunc())

	if err != nil {
		logger.Make(nil).Debug(err)

		return 0, err
	}

	// Make a Regex to say we only want letters and numbers
	regex, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		logger.Make(nil).Debug(err)

		return 0, nil
	}

	formulaVarStr := regex.ReplaceAllString(v.Formula, " ")
	formulaVarStr = strings.TrimLeft(formulaVarStr, " ")
	formulaVarStr = strings.TrimRight(formulaVarStr, " ")
	formulaVar := strings.Split(formulaVarStr, " ")

	// get formula parameters
	parameters := make(map[string]interface{}, 8)

	for _, fVar := range formulaVar {
		if contains(evalFunc, fVar) {
			continue
		}

		parameters[fVar], err = getField(v, fVar)

		if err != nil {
			parameters[fVar], err = getField(payloadValidator, fVar)
		}

		if err != nil {
			formulaVar = remove(formulaVar, fVar)
		}
	}

	if len(parameters) != len(formulaVar) {
		return 0, ErrFormulaSetup
	}

	result, err := expression.Evaluate(parameters)

	if err != nil {
		logger.Make(nil).Debug(err)

		return 0, err
	}

	// Parse interface to float
	return getFloat(result)
}

// GetVoucherResult to get voucher benefit
func (v *Validator) GetVoucherResult() (int64, error) {
	if v.ValueVoucherID == nil {
		return 0, nil
	}

	if *v.ValueVoucherID == 0 {
		return 0, nil
	}

	return *v.ValueVoucherID, nil
}

// CalculateDiscount to get the discount currency value
func (v *Validator) CalculateDiscount(trxAmount *float64) (float64, error) {
	if v.Discount == nil || *v.Discount == 0 {
		return 0, nil
	}

	result := *trxAmount * (*v.Discount / 100)

	return result, nil
}

// CalculateMaximumValue to get maximum value of a reward
func (v *Validator) CalculateMaximumValue(value *float64) (float64, error) {
	maxValue := v.MaxValue

	if maxValue == nil || *maxValue == 0 {
		return 0, nil

	}

	if *value < *maxValue {
		return *value, nil
	}

	return *maxValue, nil
}

func (v *Validator) CalculateIncentive(pv *PayloadValidator) float64 {
	var rewardResult float64
	var i Incentive

	// check if incentive is available
	if reflect.DeepEqual(i, v.Incentive) {
		return 0
	}

	// check if incentive validator is available
	if v.Incentive.Validator == nil {
		return 0
	}

	// get reward value by looping reward incentive validator that match payload validator
	for _, validator := range *v.Incentive.Validator {
		value, _ := validator.GetRewardValue(pv)
		rewardResult += value
	}

	return rewardResult
}

func getField(obj interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(obj)
	val := reflect.Indirect(r).FieldByName(strings.Title(field))

	if !val.IsValid() {
		return 0, ErrFormulaNotValid
	}

	f := val.Interface()

	if f == nil {
		return 0, nil
	}

	switch f.(type) {
	case *float64:
		v := getInterfaceValue(f)
		result, _ := strconv.ParseFloat(v, 64)

		return result, nil
	case *int64:
		v := getInterfaceValue(f)
		result, _ := strconv.ParseInt(v, 10, 64)

		return result, nil
	}

	result := getInterfaceValue(f)

	return result, nil
}

func contains(strings []string, str string) bool {
	for _, n := range strings {
		if str == n {
			return true
		}
	}

	return false
}

func getInterfaceValue(intfc interface{}) string {
	switch v := intfc.(type) {
	case *float64:
		if v == nil {
			return ""
		}

		return fmt.Sprintf("%v", *v)
	case *int64:
		if v == nil {
			return ""
		}

		return fmt.Sprintf("%v", *v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getFloat(unk interface{}) (float64, error) {
	floatType := reflect.TypeOf(float64(0))
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}

	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func goValuateFunc() map[string]govaluate.ExpressionFunction {
	return map[string]govaluate.ExpressionFunction{
		"floor": func(arguments ...interface{}) (interface{}, error) {
			return math.Floor(arguments[0].(float64)), nil
		},
	}
}
