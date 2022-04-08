package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	RequestID string                 `json:"requestID,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
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
func Make(c echo.Context, data ...map[string]interface{}) *logrus.Entry {
	var rl requestLogger
	dataMap := map[string]interface{}{}

	logrus.SetReportCaller(true)

	if c != nil {
		rl.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	if len(data) > 0 {
		dataMap = data[0]
	}

	payloadExcluder(&dataMap)
	rl.Data = dataMap

	return logrus.WithFields(logrus.Fields{
		"_requestID": rl.RequestID,
		"data":       rl.Data,
	})
}

// MakeWithoutReportCaller to get a log without report caller
func MakeWithoutReportCaller(c echo.Context, data ...map[string]interface{}) *logrus.Entry {
	log := Make(c, data...)
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

func DumpJson(args ...interface{}) {
	var b []byte
	var err error

	for idx, data := range args {
		b, err = json.Marshal(data)

		if err != nil {
			fmt.Println(err)
		}

		args[idx] = string(b)
	}

	fmt.Println("DEBUGGING ONLY")
	spew.Dump(args)
	fmt.Println("DEBUGGING ONLY")
}

func reExcludePayload(pl interface{}) (map[string]interface{}, bool) {
	vMap, ok := pl.(map[string]interface{})

	if !ok {
		return map[string]interface{}{}, ok
	}

	payloadExcluder(&vMap)

	return vMap, true
}

func payloadExcluder(pl *map[string]interface{}) {
	var ok bool
	var vMap map[string]interface{}
	plMap := *pl

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

	*pl = plMap
}

func contains(strIncluder []string, str string) bool {
	for _, include := range strIncluder {
		if strings.Contains(strings.ToLower(str), include) {
			return true
		}
	}

	return false
}
