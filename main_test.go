package main

import (
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("Ping", func() {
		config.LoadEnv()
		e := helper.NewDummyEcho("GET", "/")
		result := `{"status":"Success","message":"PONG!!","data":null}` + "\n"
		err := ping(e.Context)

		It("happy test", func() {
			Expect(err).To(BeNil())
			Expect(e.Response.Body.String()).To(Equal(result))
		})
	})
})
