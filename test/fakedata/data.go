package fakedata

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func Campaign() models.Campaign {
	metadata := CampaignMetadata()
	startDate := time.Now().AddDate(0, 0, -1)

	return models.Campaign{
		Name:        gofakeit.Name(),
		Description: gofakeit.LoremIpsumSentence(5),
		EndDate:     startDate.AddDate(1, 0, 0).Format(models.DateFormat),
		StartDate:   startDate.Format(models.DateFormat),
		Status:      helper.CreateInt8(int8(gofakeit.Number(0, 1))),
		Metadata:    &metadata,
	}
}

func CampaignMetadata() models.Metadata {
	return models.Metadata{
		IsReferral: gofakeit.Bool(),
		Prefix:     gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{5}"),
	}
}

func Reward() models.Reward {
	validator := Validator()

	return models.Reward{
		Name:               gofakeit.PetName(),
		Description:        gofakeit.LoremIpsumSentence(5),
		TermsAndConditions: gofakeit.LoremIpsumSentence(5),
		HowToUse:           gofakeit.LoremIpsumSentence(5),
		PromoCode:          gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{10}"),
		JournalAccount:     strings.ToUpper(gofakeit.DigitN(10)),
		IsPromoCode:        helper.CreateInt64(int64(gofakeit.Number(0, 1))),
		Type:               helper.CreateInt64(int64(gofakeit.Number(1, 5))),
		Validators:         &validator,
	}
}

func Validator() models.Validator {
	return models.Validator{
		Channel:         gofakeit.RandomString([]string{"9997", "9998", "9999"}),
		Product:         gofakeit.RandomString([]string{"62", "37", "36"}),
		TransactionType: gofakeit.RandomString([]string{"OP", "SL", "CC"}),
		Value:           helper.CreateFloat64(50000),
	}
}
