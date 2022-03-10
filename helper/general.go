package helper

import (
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/logger"
)

func CreateFloat64(input float64) *float64 {
	return &input
}

func CreateInt64(input int64) *int64 {
	return &input
}

func ToJson(obj interface{}) string {
	b, err := json.Marshal(obj)

	if err != nil {
		logger.Make(nil, nil).Fatal(err)
	}

	return string(b)
}

func FloatToString(input float64) string {
	return fmt.Sprintf("%f", input)
}
