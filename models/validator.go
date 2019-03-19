package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/labstack/gommon/log"
)

// Validator to store all validator data
type Validator struct {
	Channel            string   `json:"channel,omitempty"`
	Product            string   `json:"product,omitempty"`
	TransactionType    string   `json:"transactionType,omitempty"`
	Unit               string   `json:"unit,omitempty"`
	Multiplier         *float64 `json:"multiplier,omitempty"`
	Value              *int64   `json:"value,omitempty"`
	Formula            string   `json:"formula,omitempty"`
	MinimalTransaction string   `json:"minimalTransaction,omitempty"`
	Source             string   `json:"source,omitempty"` // device name that user used
}

// PayloadValidator to store a payload to validate a request
type PayloadValidator struct {
	PromoCode         string     `json:"promoCode,omitempty"`
	VoucherID         string     `json:"voucherId,omitempty"`
	CampaignID        string     `json:"campaignId,omitempty"`
	UserID            string     `json:"userId,omitempty"`
	TransactionAmount float64    `json:"transactionAmount,omitempty"`
	Validators        *Validator `json:"validators,omitempty"`
}

var skippedValidator = []string{"multiplier", "value", "formula"}
var compareEqual = []string{"channel", "product", "transactionType", "unit", "source"}

// Validate to validate client input with admin input
func (v *Validator) Validate(payloadValidator *PayloadValidator) error {
	var reqValidator map[string]interface{}

	if v == nil {
		log.Error(ErrValidatorUnavailable)

		return ErrValidatorUnavailable
	}

	vReflector := reflect.ValueOf(v).Elem()
	tempJSON, _ := json.Marshal(payloadValidator.Validators)
	json.Unmarshal(tempJSON, &reqValidator)

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
				log.Warn(ErrValidation)

				return ErrValidation
			}
		case fieldName == "minimalTransaction":
			minTrx, _ := strconv.ParseFloat(fieldValue, 64)

			if minTrx > payloadValidator.TransactionAmount {
				log.Warn(ErrValidation)

				return ErrValidation
			}
		}
	}

	return nil
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
