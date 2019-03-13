package services

// Validator to store all validator data
type Validator struct {
	Channel            string   `json:"channel,omitempty"`
	Product            string   `json:"product,omitempty"`
	TransactionType    string   `json:"transactionType,omitempty"`
	Unit               string   `json:"unit,omitempty"`
	Multiplier         *float64 `json:"multiplier,omitempty"`
	Value              *int64   `json:"value,omitempty"`
	Formula            string   `json:"formula,omitempty"`
	MinimalTransaction string   `json:"minimalTransaction,omitempty"`
}

// func Test(a interfaces.Aprinter) {
// 	// models.Cacing()
// }

// func Test1(a models.IErrors) {
// 	// models.Cacing()
// }

// Validate to validate client input with admin input
// func (v *Validator) Validate(validateVoucher PathVoucher) {
// 	var payloadValidator map[string]interface{}

// 	if v == nil {
// 		log.Error(ErrValidatorUnavailable)
// 		return nil, ErrValidatorUnavailable
// 	}

// 	vReflector := reflect.ValueOf(v).Elem()
// 	tempJSON, _ := json.Marshal(validateVoucher.Validators)
// 	json.Unmarshal(tempJSON, &payloadValidator)

// 	for i := 0; i < vReflector.NumField(); i++ {
// 		fieldName := strcase.ToLowerCamel(vReflector.Type().Field(i).Name)
// 		fieldValue := vReflector.Field(i).Interface()

// 		switch fieldName {
// 		case "channel", "product", "transactionType", "unit":
// 			if fieldValue != payloadValidator[fieldName] {
// 				log.Error(ErrValidation)
// 				return nil, ErrValidation
// 			}
// 		case "minimalTransaction":
// 			minTrx, _ := strconv.ParseFloat(fieldValue.(string), 64)

// 			if minTrx > validateVoucher.TransactionAmount {
// 				log.Error(ErrValidation)
// 				return nil, ErrValidation
// 			}
// 		}
// 	}
// }

// NewValidate to validate client input with admin input
// func (v *Validator) NewValidate() error {
// 	// return error.ErrValidation
// }

// var err interfaces.IErrors = models.Errors{}

// func cacing() error {
// 	return err.ErrValidation()
// }
