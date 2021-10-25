package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	timestampFormat = "2006-01-02 15:04:05.000"
	starString      = "**********"
)

var (
	strExclude = []string{"password", "base64", "npwp", "phone", "nik", "ktp", "gaji", "othr",
		"slik"}
)

type requestLogger struct {
	RequestID string      `json:"requestID,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

// Init function to make an initial logger
func Init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{
		TimestampFormat: timestampFormat,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			tmp := strings.Split(f.File, "/")
			filename := tmp[len(tmp)-1]
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)
}

// Make is to get a log parameter
func Make(c echo.Context, payload interface{}) *logrus.Entry {
	var rl requestLogger
	var pl []byte
	logrus.SetReportCaller(true)

	if c != nil {
		rl.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	pl, bytesOk := payload.([]byte)

	if !bytesOk {
		plTemp, err := json.Marshal(payload)

		if err != nil {
			logrus.WithFields(logrus.Fields{"requestID": rl.RequestID}).Debug(err)

			return logrus.WithFields(logrus.Fields{})
		}

		pl = plTemp
	}

	if payload != nil {
		payloadExcluder(&pl)
		rl.Payload = string(pl)
	}

	return logrus.WithFields(logrus.Fields{"params": structs.Map(rl)})
}

// MakeWithoutReportCaller to get a log without report caller
func MakeWithoutReportCaller(c echo.Context, payload interface{}) *logrus.Entry {
	log := Make(c, payload)
	logrus.SetReportCaller(false)

	return log
}

// GetEchoRID to get echo request ID
func GetEchoRID(c echo.Context) string {
	if c == nil {
		return "self-request-" + time.Now().Format("20060102150405.999")
	}

	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// MakeStructToJSON to get a json string of struct
// JUST FOR DEBUGGING TOOL
func Dump(strct ...interface{}) {
	fmt.Println("DEBUGGING ONLY")
	spew.Dump(strct)
	fmt.Println("DEBUGGING ONLY")
}

func reExcludePayload(pl interface{}) (map[string]interface{}, bool) {
	var vBytes []byte
	vMap, ok := pl.(map[string]interface{})

	if !ok {
		return map[string]interface{}{}, ok
	}

	vBytes, err := json.Marshal(vMap)

	if err != nil {
		Make(nil, nil).Error(err)

		return map[string]interface{}{}, false
	}

	payloadExcluder(&vBytes)
	err = json.Unmarshal(vBytes, &vMap)

	if err != nil {
		Make(nil, nil).Error(err)

		return map[string]interface{}{}, false
	}

	return vMap, true
}

func payloadExcluder(pl *[]byte) {
	var ok bool
	var plMap, vMap map[string]interface{}
	err := json.Unmarshal(*pl, &plMap)

	if err != nil {
		Make(nil, nil).Error(err)

		return
	}

	for k, v := range plMap {
		vMap, ok = reExcludePayload(v)

		if ok {
			plMap[k] = vMap
			continue
		}

		if contains(strExclude, k) {
			v = starString
		}

		plMap[k] = v
	}

	*pl, err = json.Marshal(plMap)

	if err != nil {
		Make(nil, nil).Error(err)

		return
	}
}

func contains(strIncluder []string, str string) bool {
	for _, include := range strIncluder {
		if strings.Contains(strings.ToLower(str), include) {
			return true
		}
	}

	return false
}
