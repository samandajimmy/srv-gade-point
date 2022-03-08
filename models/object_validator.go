package models

import (
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
	var interfaceValue interface{}
	var fieldName string
	var fieldValue string
	var rObject reflect.Value

	if validator == nil {
		return ErrValidatorUnavailable
	}

	tv := ov.TightenValidator
	vReflector := reflect.ValueOf(validator).Elem()
	mappedObj := map[string]interface{}{}

	for _, obj := range objects {
		rObject = reflect.ValueOf(obj).Elem()

		for i := 0; i < rObject.NumField(); i++ {
			interfaceValue = rObject.Field(i).Interface()
			fieldName = strcase.ToLowerCamel(rObject.Type().Field(i).Name)
			hasElement := rObject.Field(i)

			if hasElement.Kind() == reflect.Ptr && !hasElement.IsNil() {
				hasElement = hasElement.Elem()
				interfaceValue = hasElement.Interface()
			}

			if hasElement.Kind() == reflect.Struct {
				continue
			}

			mappedObj[fieldName] = interfaceValue
		}
	}

	for i := 0; i < vReflector.NumField(); i++ {
		interfaceValue = vReflector.Field(i).Interface()
		fieldName = strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
		fieldValue = getInterfaceValue(interfaceValue)

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
