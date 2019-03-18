package models

import (
	"encoding/json"
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
func (v *Validator) GetValidatorKeys(payloadValidator *PayloadValidator) error {
	return nil
}
