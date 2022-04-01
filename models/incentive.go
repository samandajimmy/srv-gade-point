package models

type DetailIncentive struct {
	PerDay        float64 `json:"perDay"`
	PerMonth      float64 `json:"perMonth"`
	ValidPerDay   bool    `json:"validPerDay"`
	ValidPerMonth bool    `json:"validPerMonth"`
	IsValid       bool    `json:"isValid"`
}

type ObjIncentive struct {
	DetailIncentive

	PerTransaction   float64 `json:"perTransaction"`
	ValidTransaction bool    `json:"validTransaction"`
}

type Incentive struct {
	MaxPerTransaction     *float64     `json:"maxPerTransaction,omitempty"`
	MaxPerDay             *float64     `json:"maxPerDay,omitempty"`
	MaxPerMonth           *float64     `json:"maxPerMonth,omitempty"`
	OslInactiveValidation bool         `json:"oslInactiveValidation,omitempty"`
	Validator             *[]Validator `json:"validator,omitempty"`
}

func (i *Incentive) Validate(obj *ObjIncentive) {
	ov := ObjectValidator{
		SkippedValidator: []string{"validator", "reward"},
		SkippedError:     []string{"maxPerDay", "maxPerMonth", "maxPerTransaction"},
		CompareEqual:     []string{},
		TightenValidator: map[string]string{
			"maxPerDay":         "perDay",
			"maxPerMonth":       "perMonth",
			"maxPerTransaction": "perTransaction",
		},
		StatusField: map[string]bool{
			"maxPerDay":         true,
			"maxPerMonth":       true,
			"maxPerTransaction": true,
		},
	}

	_ = ov.autoValidating(i, obj)
	obj.ValidPerMonth = ov.StatusField["maxPerMonth"]
	obj.ValidPerDay = ov.StatusField["maxPerDay"]
	obj.ValidTransaction = ov.StatusField["maxPerTransaction"]

	if !ov.StatusField["maxPerMonth"] {
		obj.PerMonth = *i.MaxPerMonth
	}

	if !ov.StatusField["maxPerDay"] {
		obj.PerDay = *i.MaxPerDay
	}

	if obj.ValidPerDay && obj.ValidPerMonth && obj.ValidTransaction {
		obj.IsValid = true
	}
}
