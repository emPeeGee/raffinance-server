package account

import (
	"time"

	"github.com/emPeeGee/raffinance/internal/transaction"
)

type accountResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Currency          string    `json:"currency"`
	Balance           float64   `json:"balance" gorm:"-"`
	Color             string    `json:"color"`
	Icon              string    `json:"icon"`
	TransactionCount  *int64    `json:"transactionCount" gorm:"transaction_count"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	RateWithPrevMonth *float64  `json:"rateWithPrevMonth"`
}

type accountDetailsResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Currency         string    `json:"currency"`
	Balance          float64   `json:"balance" gorm:"-"`
	TransactionCount *int64    `json:"transactionCount" gorm:"transaction_count"`
	Color            string    `json:"color"`
	Icon             string    `json:"icon"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	// Transactions     []transaction.TransactionResponse `json:"transactions" gorm:"foreignkey:to_account_id"`
	Transactions []transaction.TransactionResponse `json:"transactions" gorm:"-"`
}

type createAccountDTO struct {
	Name     string  `json:"name" validate:"required,min=2,max=256"`
	Balance  float64 `json:"balance" validate:"numeric,gte=0"`
	Currency string  `json:"currency" validate:"required,currency,min=2,max=10"`
	Icon     string  `json:"icon" validate:"required,max=128"`
	Color    string  `json:"color" validate:"required,hexcolor,min=7,max=7"`
}

type updateAccountDTO struct {
	Name     string  `json:"name" validate:"required,min=2,max=256"`
	Balance  float64 `json:"balance" validate:"numeric,gte=0"`
	Currency string  `json:"currency" validate:"required,currency,min=2,max=10"`
	Icon     string  `json:"icon" validate:"required,max=128"`
	Color    string  `json:"color" validate:"required,hexcolor,min=7,max=7"`
}
