package models

type Incentive struct {
	MaxTransaction float64     `json:"maxTransaction"`
	MaxPerDay      float64     `json:"maxPerDay"`
	MaxPerMonth    float64     `json:"maxPerMonth"`
	Validator      []Validator `json:"validator"`
}

func (i *Incentive) ValidateMaxIncentive(sumIncentive *SumIncentive) {
	ov := ObjectValidator{
		SkippedValidator: []string{"validator", "maxTransaction", "reward", "isValid"},
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

	if !sumIncentive.ValidPerDay || !sumIncentive.ValidPerMonth {
		sumIncentive.IsValid = false
	}
}

func (i *Incentive) ValidateMaxTransaction(amount float64) float64 {
	ov := ObjectValidator{
		SkippedValidator: []string{"validator", "maxPerDay", "maxPerMonth"},
		SkippedError:     []string{"maxTransaction"},
		TightenValidator: map[string]string{
			"maxTransaction": "transactionAmount",
		},
		StatusField: map[string]bool{
			"maxTransaction": true,
		},
	}

	obj := PayloadValidator{
		TransactionAmount: &amount,
	}

	_ = ov.autoValidating(i, &obj)

	if !ov.StatusField["maxTransaction"] {
		return i.MaxTransaction
	}

	return amount
}
