package validatorutil

import (
	"github.com/bojanz/currency"
	"github.com/go-playground/validator"
)

func CurrencyValidator(fl validator.FieldLevel) bool {
	currencyField := fl.Field().String()
	return currency.IsValid(currencyField)
}
