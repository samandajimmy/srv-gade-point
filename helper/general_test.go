package helper_test

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("General", func() {

	Describe("InterfaceToMap", func() {
		It("expected to return as expected", func() {
			usedCif := gofakeit.Regex("[1234567890]{10}")
			obj := models.ReqOslStatus{
				Cif:         usedCif,
				ProductCode: "02",
			}

			mappedObj := helper.InterfaceToMap(obj)
			expectedResult := map[string]interface{}{
				"cif":          usedCif,
				"product_code": "02",
			}

			Expect(mappedObj).To(Equal(expectedResult))
		})
	})
})
