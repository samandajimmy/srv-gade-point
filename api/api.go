package api

import (
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// API struct represents a request in ApiRequest
type API struct {
	Host        *url.URL
	UserAgent   string
	HTTPClient  *http.Client
	ContentType string
	Method      string
	Endpoint    string
	requestBody interface{}
	ctx         echo.Context
}

const defRTO = 30 * time.Second

var (
	// APIRCSuccess represents response code success
	APIRCSuccess = "00"

	// APIRCError represents response code error
	APIRCError = "99"
)

// NewAPI for create new client request
func NewAPI(c echo.Context, baseURL string, contentType string) (API, error) {
	url, err := url.Parse(baseURL)

	if err != nil {
		return API{}, err
	}

	return API{
		Host: url,
		HTTPClient: &http.Client{
			Timeout: defRTO,
		},
		ContentType: contentType,
		ctx:         c,
	}, nil
}

// Request represent global API Request
func (api *API) Request(pathName string, method string, body interface{}) (*http.Request, error) {
	api.requestBody = body
	api.Method = method
	api.Endpoint = pathName
	api.Host.Path += pathName
	switch ct := api.ContentType; ct {
	case echo.MIMEApplicationForm:
		return api.requestURLEncoded(method, body)
	case echo.MIMEApplicationJSON:
		return api.requestJSON(method, body)
	}

	return nil, nil
}

// Do is a function to execute the http request
func (api *API) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := api.HTTPClient.Do(req)

	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return resp, err
	}

	logger.MakeWithoutReportCaller(api.ctx, map[string]interface{}{
		"requestBody":  api.requestBody,
		"responseBody": v,
	}).Debug(fmt.Sprintf(models.ApiLogMessage, api.Host.Host, api.Endpoint))

	return resp, nil
}

func (api *API) requestURLEncoded(method string, body interface{}) (*http.Request, error) {
	var jsonData []byte
	var mapData map[string]interface{}

	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	params := url.Values{}
	err = json.Unmarshal(jsonData, &mapData)

	if err != nil {
		return nil, err
	}

	for key, value := range mapData {
		if _, ok := value.(string); !ok {
			logger.Make(nil).Debug(models.ErrSetVar)
			continue
		}

		params.Set(key, value.(string))
	}

	payloadStr, err := url.QueryUnescape(params.Encode())

	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(payloadStr)
	req, err := http.NewRequest(method, api.stringURL(), payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", echo.MIMEApplicationForm)

	return req, nil
}

func (api *API) requestJSON(method string, body interface{}) (*http.Request, error) {
	var jsonData []byte

	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(string(jsonData))
	req, err := http.NewRequest(method, api.stringURL(), payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", echo.MIMEApplicationJSON)

	return req, nil
}

func (api *API) stringURL() string {
	URLStr, _ := url.QueryUnescape(api.Host.String())

	return URLStr
}
