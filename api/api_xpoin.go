package api

import (
	"encoding/json"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/logger"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
)

const (
	XpoinCodeOslInactive = "00"
	XpoinCodeOslActive   = "404"
)

type IApiXpoin interface {
	XpoinPost(c echo.Context, body interface{}, path string) (XpoinResponse, error)
}

// APIxpoin struct represents a request for API Xpoin
type APIXpoin struct {
	Host      *url.URL
	API       API
	Method    string
	APIKey    string
	ClientId  string
	ChannelId string
}

// XpoinResponse struct represents a response for API Xpoin
type XpoinResponse struct {
	ResponseCode string                 `json:"responseCode"`
	ResponseDesc string                 `json:"responseDesc"`
	Data         string                 `json:"data"`
	ResponseData map[string]interface{} `json:"responseData,omitempty"`
}

// NewXpoinAPI is function to initiate a Xpoin API request
func NewXpoinAPI() IApiXpoin {
	apiXpoin := APIXpoin{}
	url, err := url.Parse(os.Getenv(`XPOIN_HOST`))

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return &apiXpoin
	}

	apiXpoin.Host = url
	apiXpoin.APIKey = os.Getenv(`XPOIN_API_KEY`)
	apiXpoin.ChannelId = os.Getenv(`XPOIN_CHANNEL_ID`)
	apiXpoin.ClientId = os.Getenv(`XPOIN_CLIENT_ID`)

	return &apiXpoin
}

// XpoinPost represent Post Xpoin API Request
func (xpoin *APIXpoin) XpoinPost(c echo.Context, body interface{}, path string) (XpoinResponse, error) {
	api, err := NewAPI(c, xpoin.Host.String(), echo.MIMEApplicationForm)

	if err != nil {
		return XpoinResponse{}, err
	}

	xpoin.API = api
	reqBody := helper.InterfaceToMap(body)
	reqBody["key"] = xpoin.APIKey
	reqBody["client_id"] = xpoin.ClientId
	reqBody["channel_id"] = xpoin.ChannelId
	req, err := xpoin.request(path, echo.POST, reqBody)

	if err != nil {
		return XpoinResponse{}, err
	}

	resp := XpoinResponse{}
	_, err = xpoin.do(req, &resp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return XpoinResponse{}, err
	}

	return resp, nil

}

// Do is a function to execute the http request
func (xpoin *APIXpoin) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := xpoin.API.Do(req, v)

	if err != nil {
		return nil, err
	}

	err = xpoin.mappingDataResponseXpoin(v)

	if err != nil {
		logger.Make(xpoin.API.ctx, nil).Debug(err)

		return resp, err
	}

	return resp, err
}

// Request represent global API Request
func (xpoin *APIXpoin) request(endpoint, method string, body interface{}) (*http.Request, error) {
	xpoin.Method = method
	req, err := xpoin.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (xpoin *APIXpoin) mappingDataResponseXpoin(v interface{}) error {
	resp := v.(*XpoinResponse)

	if resp.Data == "" {
		return nil
	}

	err := json.Unmarshal([]byte(resp.Data), &resp.ResponseData)

	if err != nil {
		logger.Make(xpoin.API.ctx, nil).Fatal("Response Data Error Unmarshal")
	}

	resp.Data = ""

	return nil
}
