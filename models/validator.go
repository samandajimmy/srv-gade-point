package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

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
}

// PayloadValidator to store a payload to validate a request
type PayloadValidator struct {
	PromoCode         string     `json:"promoCode,omitempty"`
	VoucherID         string     `json:"voucherId,omitempty"`
	UserID            string     `json:"userId,omitempty"`
	TransactionAmount float64    `json:"transactionAmount,omitempty"`
	Validators        *Validator `json:"validators,omitempty"`
}

var skippedValidator = []string{"userID", "transactionAmount", "multiplier", "value", "formaula"}

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
		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue := vReflector.Field(i).Interface()

		switch fieldName {
		case "channel", "product", "transactionType", "unit":
			if fieldValue != reqValidator[fieldName] {
				log.Error(ErrValidation)

				return ErrValidation
			}
		case "minimalTransaction":
			minTrx, _ := strconv.ParseFloat(fieldValue.(string), 64)

			if minTrx > payloadValidator.TransactionAmount {
				log.Error(ErrValidation)

				return ErrValidation
			}
		}
	}

	return nil
}

// GetValidatorKeys to get all validator keys needed
func (v *Validator) GetValidatorKeys(payloadValidator *GetCampaignValue) map[string]string {
	validator := make(map[string]string)
	vReflector := reflect.ValueOf(payloadValidator).Elem()
	var value string

	for i := 0; i < vReflector.NumField(); i++ {
		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue := vReflector.Field(i).Interface()

		if contains(skippedValidator, fieldName) {
			continue
		}

		switch fieldValue.(type) {
		case float64:
			value = fmt.Sprintf("%f", fieldValue.(float64))

			validator[fieldName] = value
		default:
			value, ok := fieldValue.(string)

			if !ok {
				log.Error(ok)
			}

			validator[fieldName] = value
		}

	}

	return validator
}

func contains(strings []string, str string) bool {
	for _, n := range strings {
		if str == n {
			return true
		}
	}

	return false
}
