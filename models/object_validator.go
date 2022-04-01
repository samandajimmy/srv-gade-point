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

		mappedObj = ov.mappingObj(rObject, mappedObj)
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

func (ov *ObjectValidator) mappingObj(rObject reflect.Value, mappedObj map[string]interface{}) map[string]interface{} {
	var hasElement reflect.Value
	var interfaceValue interface{}
	var fieldName string

	obj := map[string]interface{}{}

	for i := 0; i < rObject.NumField(); i++ {
		hasElement = rObject.Field(i)
		interfaceValue = hasElement.Interface()
		fieldName = strcase.ToLowerCamel(rObject.Type().Field(i).Name)

		if contains(ov.SkippedValidator, fieldName) {
			continue
		}

		if hasElement.Kind() == reflect.Ptr && !hasElement.IsNil() {
			hasElement = hasElement.Elem()
			interfaceValue = hasElement.Interface()
		}

		if hasElement.Kind() == reflect.Struct {
			obj = ov.mappingObj(hasElement, obj)
			continue
		}

		obj[fieldName] = interfaceValue
	}

	for k, v := range obj {
		mappedObj[k] = v
	}

	return mappedObj
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
