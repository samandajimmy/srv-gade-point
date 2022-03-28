package fakedata

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/viper"
)

func Campaign() models.Campaign {
	startDate := time.Now().AddDate(0, 0, -1)

	return models.Campaign{
		Name:        gofakeit.Name(),
		Description: gofakeit.LoremIpsumSentence(5),
		EndDate:     startDate.AddDate(1, 0, 0).Format(models.DateFormat),
		StartDate:   startDate.Format(models.DateFormat),
		Status:      helper.CreateInt8(1),
	}
}

func CampaignReferral() models.Campaign {
	campaign := Campaign()
	metadata := CampaignMetadata()
	metadata.IsReferral = true
	campaign.Metadata = &metadata

	return campaign
}

func VoucherDirectDisc() models.Voucher {
	startDate := time.Now().AddDate(0, 0, -1)
	stock := int32(5)
	validator := Validator()
	value := gofakeit.Float64Range(50000, 1000000)
	validator.Value = &value

	return models.Voucher{
		Name:               gofakeit.PetName(),
		Description:        gofakeit.LoremIpsumSentence(10),
		StartDate:          startDate.Format(models.DateFormat),
		EndDate:            startDate.AddDate(1, 0, 0).Format(models.DateFormat),
		Point:              helper.CreateInt64(int64(gofakeit.Number(1, 10))),
		JournalAccount:     strings.ToUpper(gofakeit.DigitN(10)),
		Type:               &models.VoucherTypeDirectDiscount,
		ImageURL:           gofakeit.URL(),
		Status:             helper.CreateInt8(int8(models.CampaignActive)),
		GeneratorType:      helper.CreateInt8(0),
		Stock:              &stock,
		PrefixPromoCode:    gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{5}"),
		TermsAndConditions: gofakeit.LoremIpsumSentence(10),
		HowToUse:           gofakeit.LoremIpsumSentence(10),
		Validators:         &validator,
		Synced:             true,
	}
}

func CampaignMetadata() models.Metadata {
	return models.Metadata{
		IsReferral: gofakeit.Bool(),
		Prefix:     gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{5}"),
	}
}

func Reward(withPromoCode bool) models.Reward {
	reward := models.Reward{
		Name:               gofakeit.PetName(),
		Description:        gofakeit.LoremIpsumSentence(5),
		TermsAndConditions: gofakeit.LoremIpsumSentence(5),
		HowToUse:           gofakeit.LoremIpsumSentence(5),
		JournalAccount:     strings.ToUpper(gofakeit.DigitN(10)),
		IsPromoCode:        helper.CreateInt64(0),
		Type:               helper.CreateInt64(int64(gofakeit.Number(1, 5))),
	}

	if withPromoCode {
		reward.IsPromoCode = helper.CreateInt64(1)
		reward.PromoCode = gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{10}")
	}

	return reward
}

func Validator() models.Validator {
	return models.Validator{
		Channel:         gofakeit.RandomString([]string{"9997", "9998", "9999"}),
		Product:         gofakeit.RandomString([]string{"62", "37", "36"}),
		TransactionType: gofakeit.RandomString([]string{"OP", "SL", "CC"}),
	}
}

func Incentive(validators ...models.Validator) models.Incentive {
	incentive := models.Incentive{
		MaxTransaction: helper.CreateFloat64(gofakeit.Float64Range(50000, 1000000)),
		MaxPerDay:      helper.CreateFloat64(gofakeit.Float64Range(50000, 1000000)),
		MaxPerMonth:    helper.CreateFloat64(gofakeit.Float64Range(50000, 1000000)),
		Validator:      &validators,
	}

	return incentive
}

func RewardDirectDisc(withPromoCode bool) models.Reward {
	reward := Reward(withPromoCode)
	reward.Type = &models.RewardTypeDirectDiscount

	reward.Validators = func() *models.Validator {
		validator := Validator()
		value := gofakeit.Float64Range(50000, 1000000)
		validator.Value = &value

		return &validator
	}()

	return reward
}

func RewardPercentDisc(withPromoCode, withMax bool) models.Reward {
	reward := Reward(withPromoCode)
	reward.Type = &models.RewardTypePercentageDiscount
	reward.Validators = func() *models.Validator {
		validator := Validator()
		disc := gofakeit.Float64Range(1, 100)
		maxValue := gofakeit.Float64Range(50000, 1000000)
		validator.Discount = &disc

		if withMax {
			validator.MaxValue = &maxValue
		}

		return &validator
	}()

	return reward
}

func RewardGoldback(withPromoCode bool) models.Reward {
	reward := Reward(withPromoCode)
	reward.Type = &models.RewardTypeGoldback
	reward.Validators = func() *models.Validator {
		validator := Validator()
		validator.IsDecimal = true
		validator.Formula = "transactionAmount/deviden"
		deviden := gofakeit.Float64Range(50000, 1000000)
		validator.Deviden = &deviden

		return &validator
	}()

	return reward
}

func RewardVoucher(withPromoCode bool, voucherId int64) models.Reward {
	reward := Reward(withPromoCode)
	reward.Type = &models.RewardTypeVoucher
	reward.Validators = func() *models.Validator {
		validator := Validator()
		validator.ValueVoucherID = &voucherId

		return &validator
	}()

	return reward
}

func RewardIncentive(withPromoCode bool) models.Reward {
	var validator models.Validator
	_ = viper.UnmarshalKey("validator.with_incentive", &validator)
	reward := Reward(withPromoCode)
	reward.Type = &models.RewardTypeIncentive
	reward.Validators = &validator

	return reward
}

func PayloadInquiry() models.PayloadValidator {
	validator := Validator()

	return models.PayloadValidator{
		BranchCode:        gofakeit.Regex("[1234567890]{5}"),
		CIF:               gofakeit.Regex("[1234567890]{10}"),
		CustomerName:      gofakeit.PetName(),
		Phone:             "0815" + gofakeit.Regex("[1234567890]{10}"),
		PromoCode:         gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{10}"),
		TransactionDate:   time.Now().Format(models.DateTimeFormat + ".000"),
		TransactionAmount: helper.CreateFloat64(gofakeit.Float64Range(50000, 1000000)),
		Validators:        &validator,
	}
}
