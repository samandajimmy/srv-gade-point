package http

import (
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"net/http"

	"github.com/labstack/echo"
)

var response models.Response

// RewardHandler represent the httphandler for reward
type RewardHandler struct {
	RewardUseCase rewards.UseCase
}

// NewRewardHandler represent to register reward endpoint
func NewRewardHandler(echoGroup models.EchoGroup, us rewards.UseCase) {
	handler := &RewardHandler{
		RewardUseCase: us,
	}

	// Fake End Point
	echoGroup.API.POST("/rewards/inquiry", handler.rewardInquiry)
	echoGroup.API.POST("/rewards/succeeded", handler.rewardSucceeded)
}

func (rwd *RewardHandler) rewardInquiry(echTx echo.Context) error {
	var reward models.Reward
	response = models.Response{}
	err := echTx.Bind(&reward)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	responseData := []byte(`{"status":"Success", "message":"Data Successfully Sent", "data":{"refTrx":"1122334455","rewards":[{"discountAmount":"20%","journalAccount":"1122334455","type":"discount","unit":"rupiah","value":200000},{"journalAccount":"1122334455","type":"discount","unit":"rupiah","value":100000},{"journalAccount":"1122334455","type":"goldback","unit":"gram","value":10},{"journalAccount":"1122334455","type":"voucher","promoCode":"PSN0JLK8","voucherName":"Voucher diskon 1.000 top up tabungan emas"}]}}`)
	var rawData map[string]interface{}
	json.Unmarshal(responseData, &rawData)
	data, _ := json.Marshal(rawData)

	return echTx.JSONBlob(http.StatusOK, data)
}

func (rwd *RewardHandler) rewardSucceeded(echTx echo.Context) error {
	var reward models.Reward
	response = models.Response{}
	err := echTx.Bind(&reward)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess

	return echTx.JSON(http.StatusOK, response)
}
