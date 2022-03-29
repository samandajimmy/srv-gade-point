package usecase

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

func (rwd *rewardUseCase) UPostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) ([]models.CoreTrxResponse, error) {

	if len(coreTrx) == 0 {
		return nil, models.ErrCoreTrxEmpty
	}

	for i, inq := range coreTrx {

		/* Flow
		* 1. Call usecase Inquiry
		* 2. Call usecae Payment Success
		 */
		var validator models.Validator
		validator.TransactionType = inq.TrxType
		validator.Product = inq.ProductCode
		validator.Channel = inq.Channel

		var plInq models.PayloadValidator
		plInq.Validators = &validator
		plInq.CIF = inq.CIF
		plInq.Phone = inq.PhoneNumber
		plInq.TransactionAmount = &inq.TrxAmount
		plInq.PromoCode = inq.MarketingCode
		plInq.TransactionDate = inq.TrxDate
		plInq.CustomerName = inq.CustomerName
		plInq.IsMulti = true

		// loan ammount jika ada maka harus di isi
		if inq.LoanAmount > 0 {
			plInq.LoanAmount = &inq.LoanAmount
		}

		// intesrest ammount jika ada maka harus di isi
		if inq.InterestAmount > 0 {
			plInq.InterestAmount = &inq.InterestAmount
		}

		var resInq models.RewardsInquiry
		resInq, respErrors := rwd.Inquiry(c, &plInq)

		if respErrors.Title != "" {
			coreTrx[i].InqStatus = models.CoreTrxInqFailed
			continue
		}

		coreTrx[i].InqStatus = URwdPayment(c, rwd, resInq, inq)
		coreTrx[i].RefTrx = resInq.RefTrx
	}

	var responseData []models.CoreTrxResponse
	err := rwd.rewardRepo.RPostCoreTrx(c, coreTrx)

	if err != nil {
		return nil, models.ErrCoreTrxFailed
	}

	return responseData, nil
}

// use in reward payment after inqury success
func URwdPayment(c echo.Context, rwd *rewardUseCase, resInq models.RewardsInquiry, inq models.CoreTrxPayload) string {
	var plPayment models.RewardPayment
	var status string = models.CoreTrxSuccess

	for _, rwds := range *resInq.Rewards {
		plPayment.RefCore = inq.TrxID
		plPayment.CIF = rwds.CIF
		plPayment.RefTrx = rwds.RefTrx

		_, errResp := rwd.Payment(c, &plPayment)

		if errResp != nil {
			status = models.CoreTrxPaymentFailed
			break
		}

	}

	return status
}
