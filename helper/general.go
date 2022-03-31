package helper

import (
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/logger"
	"math/rand"
	"time"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func CreateFloat64(input float64) *float64 {
	return &input
}

func CreateInt64(input int64) *int64 {
	return &input
}

func CreateInt8(input int8) *int8 {
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

func RandomStr(n int, arrChecker map[string]bool) string {
	var randString string
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	randString = string(b)

	if arrChecker[randString] {
		randString = string(RandomStr(n, arrChecker))
	}

	arrChecker[randString] = true

	return randString
}

func NowDbBun() time.Time {
	return time.Now().Add(7 * time.Hour)
}

func InterfaceToMap(obj interface{}) map[string]interface{} {
	var mappedObj map[string]interface{}

	byteObj, err := json.Marshal(obj)

	if err != nil {
		logger.Make(nil, nil).Error(err)
	}

	if err := json.Unmarshal(byteObj, &mappedObj); err != nil {
		logger.Make(nil, nil).Error(err)
	}

	return mappedObj
}
