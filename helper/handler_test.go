package helper_test

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var iHandler helper.IHandler
	var stHandler helper.Handler
	var e helper.DummyEcho
	var reqpl map[string]interface{}
	var pl interface{}
	var err error

	JustBeforeEach(func() {
		e = helper.NewDummyEcho(http.MethodPost, "/", reqpl)
		iHandler = helper.NewHandler(&stHandler)
	})

	BeforeEach(func() {
		stHandler = helper.Handler{}
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = iHandler.Validate(e.Context, pl)
		})

		BeforeEach(func() {
			reqpl = map[string]interface{}{}
			pl = nil
		})

		Context("error echo binding", func() {
			BeforeEach(func() {
				pl = &[]struct{ Field string }{}
			})

			It("expect to return error", func() {
				Expect(err.Error()).To(Equal("code=400, message=Unmarshal type error: expected=[]struct { Field string }, got=object, field=, offset=1"))
			})
		})

		Context("error echo validate", func() {
			BeforeEach(func() {
				reqpl = nil
				pl = models.AccountToken{}
			})

			It("expect to return error", func() {
				Expect(err).To(Equal(models.ErrInternalServerError))
			})
		})

		Context("succeeded", func() {
			It("expect to return nil error", func() {
				Expect(err).To(BeNil())
			})

			It("expect to reset struct handler", func() {
				Expect(stHandler).To(Equal(helper.Handler{}))
			})
		})
	})

	Describe("ShowResponse", func() {
		var respData models.AccountToken
		var mockResp models.Response
		var errInput error
		var errors models.ResponseErrors

		JustBeforeEach(func() {
			err = iHandler.ShowResponse(e.Context, respData, errInput, errors)
		})

		BeforeEach(func() {
			respData = models.AccountToken{}
			mockResp = models.Response{}
			errInput = nil
		})

		Context("error input is nil", func() {
			Context("response errors is nil", func() {
				BeforeEach(func() {
					respData = models.AccountToken{
						Username: "jimmy",
					}
					mockResp = models.Response{
						Code:    "00",
						Status:  "Success",
						Message: "Data Berhasil Dikirim",
						Data:    respData,
					}
				})

				It("expect to have response data ", func() {
					Expect(stHandler.Response).To(Equal(mockResp))
				})

				It("expect to have nil error ", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("response errors is not nil", func() {
				BeforeEach(func() {
					errors = models.ResponseErrors{Title: "This is errors"}
					mockResp = models.Response{
						Code:    "99",
						Status:  "Error",
						Message: "This is errors",
					}

				})

				It("expect to show response errors ", func() {
					Expect(stHandler.Response).To(Equal(mockResp))
				})

				It("expect to have nil error ", func() {
					Expect(err).To(BeNil())
				})
			})

		})

		Context("error input is not nil", func() {
			BeforeEach(func() {
				errInput = models.ErrAddQuotaFailed
				mockResp = models.Response{
					Code:    "99",
					Status:  "Error",
					Message: errInput.Error(),
				}
			})

			It("expect to have response error", func() {
				Expect(stHandler.Response).To(Equal(mockResp))
			})

			It("expect to return nil error", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("SetTotalCount", func() {
		JustBeforeEach(func() {
			iHandler.SetTotalCount("100")
		})

		It("expect to have totalCount on response total count", func() {
			Expect(stHandler.Response.TotalCount).To(Equal("100"))
		})
	})
})
