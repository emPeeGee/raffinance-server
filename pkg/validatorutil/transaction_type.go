package validatorutil

import (
	"github.com/emPeeGee/raffinance/pkg/util"
	"github.com/go-playground/validator"
)

func TransactionType(fl validator.FieldLevel) bool {
	transactionTypeField := fl.Field().Uint()
	allowedTypes := []int{1, 2, 3}

	return util.Contains(allowedTypes, int(transactionTypeField))
}
