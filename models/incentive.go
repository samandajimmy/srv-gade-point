package models

type Incentive struct {
	MaxTransaction float64     `json:"maxTransaction"`
	MaxPerDay      float64     `json:"maxPerDay"`
	MaxPerMonth    float64     `json:"maxPerMonth"`
	Validator      []Validator `json:"validator"`
}

func (i *Incentive) ValidateMaxIncentive(sumIncentive *SumIncentive) {
	ov := ObjectValidator{
		SkippedValidator: []string{"validator", "maxTransaction"},
		SkippedError:     []string{"maxPerDay", "maxPerMonth"},
		CompareEqual:     []string{},
		TightenValidator: map[string]string{
			"maxPerDay":   "perDay",
			"maxPerMonth": "perMonth",
		},
		StatusField: map[string]bool{
			"maxPerDay":   true,
			"maxPerMonth": true,
		},
	}

	_ = ov.autoValidating(i, sumIncentive)
	sumIncentive.ValidPerMonth = ov.StatusField["maxPerMonth"]
	sumIncentive.ValidPerDay = ov.StatusField["maxPerDay"]

	if !ov.StatusField["maxPerMonth"] {
		sumIncentive.PerMonth = i.MaxPerMonth
	}

	if !ov.StatusField["maxPerDay"] {
		sumIncentive.PerDay = i.MaxPerDay
	}
}

func (i *Incentive) ValidateMaxTransaction(amount float64) float64 {
	if amount > i.MaxTransaction {
		amount = i.MaxTransaction
	}

	return amount
}
