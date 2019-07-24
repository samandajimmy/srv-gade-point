package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	govaluate "gopkg.in/Knetic/govaluate.v2"
)

// Validator to store all validator data
type Validator struct {
	CampaignCode         string   `json:"campaignCode,omitempty"`
	Channel              string   `json:"channel,omitempty" validate:"required"`
	Discount             *float64 `json:"discount,omitempty"`
	Formula              string   `json:"formula,omitempty"`
	MaxLoanAmount        *float64 `json:"maxLoanAmount,omitempty"`
	MaxValue             *float64 `json:"maxValue,omitempty"`
	MinLoanAmount        *float64 `json:"minLoanAmount,omitempty"`
	MinTransactionAmount *float64 `json:"minTransactionAmount,omitempty"`
	Multiplier           *float64 `json:"multiplier,omitempty"`
	Product              string   `json:"product,omitempty" validate:"required"`
	Source               string   `json:"source,omitempty"` // device name that user used
	TransactionType      string   `json:"transactionType,omitempty" validate:"required"`
	Unit                 string   `json:"unit,omitempty"`
	Value                *float64 `json:"value,omitempty"`
	ValueVoucherID       *int64   `json:"valueVoucherId,omitempty"`
}

// PayloadValidator to store a payload to validate a request
type PayloadValidator struct {
	BranchCode        string     `json:"branchCode,omitempty"`
	CampaignID        string     `json:"campaignId,omitempty"`
	CIF               string     `json:"cif,omitempty"`
	CustomerName      string     `json:"customerName,omitempty"`
	LoanAmount        *float64   `json:"loanAmount,omitempty"`
	Phone             string     `json:"phone,omitempty"`
	PromoCode         string     `json:"promoCode,omitempty" validate:"required"`
	RedeemedDate      string     `json:"redeemedDate,omitempty"`
	RefTrx            string     `json:"refTrx,omitempty"`
	TransactionDate   string     `json:"transactionDate,omitempty" validate:"required"`
	TransactionAmount *float64   `json:"transactionAmount,omitempty" validate:"required"`
	VoucherID         string     `json:"voucherId,omitempty"`
	Validators        *Validator `json:"validators,omitempty"`
}

var skippedValidator = []string{"multiplier", "value", "formula", "maxValue", "unit"}
var compareEqual = []string{"channel", "product", "transactionType", "source", "campaignCode"}
var tightenValidator = map[string]string{
	"minTransactionAmount": "transactionAmount",
	"minLoanAmount":        "loanAmount",
	"maxLoanAmount":        "loanAmount",
}
var customErrMsg = "%s on this transaction is not valid to use the benefit"

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
	var reqValidator map[string]interface{}
	var payloadVal map[string]interface{}

	if v == nil {
		logrus.Debug(ErrValidatorUnavailable)

		return ErrValidatorUnavailable
	}

	vReflector := reflect.ValueOf(v).Elem()
	tempJSON, _ := json.Marshal(payloadValidator.Validators)
	json.Unmarshal(tempJSON, &reqValidator)
	tempJSON, _ = json.Marshal(payloadValidator)
	json.Unmarshal(tempJSON, &payloadVal)

	for i := 0; i < vReflector.NumField(); i++ {
		interfaceValue := vReflector.Field(i).Interface()

		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue := getInterfaceValue(interfaceValue)

		if contains(skippedValidator, fieldName) {
			continue
		}

		if fieldValue == "" {
			continue
		}

		switch {
		case contains(compareEqual, fieldName):
			reqValidatorVal := fmt.Sprintf("%v", reqValidator[fieldName])

			if !strings.Contains(fieldValue, reqValidatorVal) {
				logrus.Debug(fmt.Errorf(customErrMsg, fieldName))

				return fmt.Errorf(customErrMsg, fieldName)
			}
		case strings.Contains(fieldName, "min"):
			minTrx, _ := strconv.ParseFloat(fieldValue, 64)

			if minTrx > payloadVal[tightenValidator[fieldName]].(float64) {
				logrus.Debug(fmt.Errorf(customErrMsg, fieldName))

				return fmt.Errorf(customErrMsg, tightenValidator[fieldName])
			}
		case strings.Contains(fieldName, "max"):
			maxTrx, _ := strconv.ParseFloat(fieldValue, 64)

			if maxTrx < payloadVal[tightenValidator[fieldName]].(float64) {
				logrus.Debug(fmt.Errorf(customErrMsg, fieldName))

				return fmt.Errorf(customErrMsg, tightenValidator[fieldName])
			}
		}
	}

	return nil
}

// GetFormulaResult to proccess the formula then get the result
func (v *Validator) GetFormulaResult(payloadValidator *PayloadValidator) (float64, error) {
	// check formula availability
	if v.Formula == "" {
		return float64(0), nil
	}

	expression, err := govaluate.NewEvaluableExpression(v.Formula)

	if err != nil {
		logrus.Debug(err)

		return 0, err
	}

	// Make a Regex to say we only want letters and numbers
	regex, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		logrus.Debug(err)

		return 0, nil
	}

	formulaVarStr := regex.ReplaceAllString(v.Formula, " ")
	formulaVarStr = strings.TrimLeft(formulaVarStr, " ")
	formulaVarStr = strings.TrimRight(formulaVarStr, " ")
	formulaVar := strings.Split(formulaVarStr, " ")
	remove(formulaVar, "transactionAmount")

	// get formula parameters
	parameters := make(map[string]interface{}, 8)
	parameters["transactionAmount"] = *payloadValidator.TransactionAmount

	for _, fVar := range formulaVar {
		parameters[fVar], err = v.getField(fVar)

		if err != nil {
			break
		}
	}

	if err != nil {
		return 0, err
	}

	result, err := expression.Evaluate(parameters)

	if err != nil {
		logrus.Debug(err)

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

func (v *Validator) getField(field string) (interface{}, error) {
	r := reflect.ValueOf(v)
	val := reflect.Indirect(r).FieldByName(strings.Title(field))

	if !val.IsValid() {
		return 0, errors.New("Formula is not valid")
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
