package models

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo"
)

var (
	responseSuccess = "Success"
	responseError   = "Error"
)

var responseCode = map[string]string{
	responseSuccess: "00",
	responseError:   "99",
}

// Response struct is represent a data for output payload
type Response struct {
	Code        string      `json:"code,omitempty"`
	Status      string      `json:"status,omitempty"`
	Message     string      `json:"message,omitempty"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	TotalCount  string      `json:"totalCount,omitempty"`
}

// SetResponse is a function to set response
func (resp *Response) SetResponse(respData interface{}, respErrors *ResponseErrors) {
	typeResp := reflect.TypeOf(respData)

	if typeResp.Kind() != reflect.Slice && respData != reflect.Zero(typeResp).Interface() {
		resp.Data = respData
	}

	if respErrors.Title == "" {
		resp.Status = responseSuccess
		resp.Code = responseCode[responseSuccess]
		resp.Message = MessageDataSuccess
		resp.Data = respData

		return
	}

	resp.Status = responseError
	resp.Code = responseCode[responseError]
	resp.Message = respErrors.Title

	if len(respErrors.Details) != 0 {
		resp.Description = strings.Join(respErrors.Details, ", ")
	}

	if respErrors.Code != "" {
		resp.Code = respErrors.Code
	}
}

// Body is function to get response body
func (resp *Response) Body(c echo.Context, err error) error {
	return c.JSON(getStatusCode(err), resp)
}

// ResponseErrors struct is represent a data for output payload
type ResponseErrors struct {
	Title   string
	Code    string
	Details []string
}

// NewResponse is a function to do initial response
func NewResponse() (Response, ResponseErrors) {
	return Response{}, ResponseErrors{}
}

// SetTitle title of Response errors
func (re *ResponseErrors) SetTitle(title string) {
	re.Title = title
}

// SetTitleCode title and code of Response errors
func (re *ResponseErrors) SetTitleCode(code string, title string, desc string) {
	re.Title = title
	re.Code = code
	re.AddError(desc)
}

// AddError adding error on Response errors
func (re *ResponseErrors) AddError(errString string) {
	re.Details = append(re.Details, errString)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusOK
	}
}
