package transaction

import (
	"github.com/go-playground/validator"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func validateTransactionType(txnType TransactionType, fromAccountID *uint, toAccountID uint) (errs []ValidationError) {
	switch txnType {
	case EXPENSE, INCOME:
		{
			if fromAccountID != nil {
				errs = append(errs, ValidationError{
					Field: "fromAccount",
					Error: "from account is not allowed for this transaction type",
				})
			}
		}
	case TRANSFER:
		{
			// accounts mustn't be the same
			if toAccountID == *fromAccountID {
				errs = append(errs, ValidationError{
					Field: "fromAccount",
					Error: "from account and to account must be different",
				})
			}

			// when transfer, from is mult
			if fromAccountID == nil {
				errs = append(errs, ValidationError{
					Field: "fromAccount",
					Error: "from account is required for a transfer transaction",
				})
			}
		}
	}

	return errs
}

func ValidateCreateTransaction(sl validator.StructLevel) {
	txn := sl.Current().Interface().(CreateTransactionDTO)
	errs := validateTransactionType(TransactionType(txn.TransactionTypeID), txn.FromAccountID, txn.ToAccountID)

	for _, err := range errs {
		sl.ReportError(txn, err.Field, "FromAccountID", err.Error, "")
	}
}

func ValidateUpdateTransaction(sl validator.StructLevel) {
	txn := sl.Current().Interface().(UpdateTransactionDTO)
	errs := validateTransactionType(TransactionType(txn.TransactionTypeID), txn.FromAccountID, txn.ToAccountID)

	for _, err := range errs {
		sl.ReportError(txn, err.Field, "FromAccountID", err.Error, "")
	}
}
