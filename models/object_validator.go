package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type ObjectValidator struct {
	SkippedValidator []string
	SkippedError     []string
	CompareEqual     []string
	TightenValidator map[string]string
	StatusField      map[string]bool
}

func (ov *ObjectValidator) autoValidating(validator interface{}, objects ...interface{}) error {
	var tempJSON []byte

	if validator == nil {
		return ErrValidatorUnavailable
	}

	tv := ov.TightenValidator
	vReflector := reflect.ValueOf(validator).Elem()
	mappedObj := map[string]interface{}{}
	tmpMappedObj := map[string]interface{}{}

	for _, obj := range objects {
		tmpMappedObj = map[string]interface{}{}

		tempJSON, _ = json.Marshal(obj)
		_ = json.Unmarshal(tempJSON, &tmpMappedObj)

		for key, value := range tmpMappedObj {
			mappedObj[key] = value
		}
	}

	for i := 0; i < vReflector.NumField(); i++ {
		interfaceValue := vReflector.Field(i).Interface()

		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue := getInterfaceValue(interfaceValue)

		if contains(ov.SkippedValidator, fieldName) {
			continue
		}

		if fieldValue == "" {
			continue
		}

		switch {
		case contains(ov.CompareEqual, fieldName):
			reqValidatorVal := fmt.Sprintf("%v", mappedObj[fieldName])

			if ov.isError(!strings.Contains(fieldValue, reqValidatorVal), fieldName) {
				return fmt.Errorf(customErrMsg, fieldName)
			}
		case strings.Contains(fieldName, "min"):
			minTrx, _ := strconv.ParseFloat(fieldValue, 64)

			if ov.isError(minTrx > mappedObj[tv[fieldName]].(float64), fieldName) {
				return fmt.Errorf(customErrMsg, tv[fieldName])
			}
		case strings.Contains(fieldName, "max"):
			maxTrx, _ := strconv.ParseFloat(fieldValue, 64)

			if ov.isError(maxTrx < mappedObj[tv[fieldName]].(float64), fieldName) {
				return fmt.Errorf(customErrMsg, tv[fieldName])
			}
		}
	}

	return nil
}

func (ov *ObjectValidator) isError(stateError bool, fieldName string) bool {
	if !stateError {
		return false
	}

	if _, ok := ov.StatusField[fieldName]; ok {
		ov.StatusField[fieldName] = false
	}

	return !contains(ov.SkippedError, fieldName)
}
