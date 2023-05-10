package transaction

import (
	"time"

	"github.com/emPeeGee/raffinance/internal/category"
	"github.com/emPeeGee/raffinance/internal/tag"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `json:"name"`
}

type TransactionResponse struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Location    string    `json:"location"`

	FromAccountID     *uint                          `json:"fromAccountId,omitempty"`
	ToAccountID       uint                           `json:"toAccountId"`
	TransactionTypeID byte                           `json:"transactionTypeId"`
	Category          category.CategoryShortResponse `json:"category"`
	Tags              []tag.TagShortResponse         `json:"tags"`
}

type CreateTransactionDTO struct {
	Date        time.Time `json:"date" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description" validate:"omitempty,max=256"`
	Location    string    `json:"location" validate:"omitempty,max=128"`

	CategoryID uint `json:"categoryId" validate:"required,numeric"`
	// NOTE: valid order matters, unique can't be the last
	TagIDs []uint `json:"tagIds" validate:"omitempty,unique,dive,numeric,gt=0"`

	// TODO: Sending as string breaks the app
	FromAccountID     *uint `json:"fromAccountId" validate:"omitempty,numeric"`
	ToAccountID       uint  `json:"toAccountId" validate:"required,numeric"`
	TransactionTypeID byte  `json:"transactionTypeId" validate:"numeric,transactiontype"`
}

type UpdateTransactionDTO struct {
	Date        time.Time `json:"date" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description" validate:"omitempty"`
	Location    string    `json:"location" validate:"omitempty,max=128"`

	CategoryID uint `json:"categoryId" validate:"required,numeric"`
	// NOTE: valid order matters, unique can't be the last
	TagIDs []uint `json:"tagIds" validate:"omitempty,unique,dive,numeric,gt=0"`

	FromAccountID     *uint `json:"fromAccountId" validate:"omitempty,numeric"`
	ToAccountID       uint  `json:"toAccountId" validate:"required,numeric"`
	TransactionTypeID byte  `json:"transactionTypeId" validate:"numeric,transactiontype"`
}

type TransactionFilter struct {
	userID      *uint
	Type        *byte
	StartDate   *time.Time
	EndDate     *time.Time
	Day         *time.Time
	Accounts    []uint
	Categories  []uint
	Tags        []uint
	Description string
}
