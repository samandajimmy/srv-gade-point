package helper

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

type IHandler interface {
	Validate(c echo.Context, pl interface{}) error
	SetTotalCount(counter string)
	ShowResponse(c echo.Context, respData interface{}, err error, errors models.ResponseErrors) error
}

type Handler struct {
	Response   models.Response
	RespErrors models.ResponseErrors
}

func NewHandler(h *Handler) IHandler {
	return h
}

func (h *Handler) Validate(c echo.Context, pl interface{}) error {
	h.reset()

	if err := c.Bind(&pl); err != nil {
		h.setError(err)

		return err
	}

	if err := c.Validate(pl); err != nil {
		h.setError(err)

		return err
	}

	return nil
}

func (h *Handler) ShowResponse(c echo.Context, respData interface{}, err error, errors models.ResponseErrors) error {
	if err != nil {
		h.setError(err)

		return h.Response.Body(c, err)
	}

	h.Response.SetResponse(respData, &errors)

	return h.Response.Body(c, err)
}

func (h *Handler) SetTotalCount(counter string) {
	h.Response.TotalCount = counter
}

func (h *Handler) setError(err error) {
	h.RespErrors.SetTitle(err.Error())
	h.Response.SetResponse("", &h.RespErrors)
}

func (h *Handler) reset() {
	h.Response = models.Response{}
	h.RespErrors = models.ResponseErrors{}
}
