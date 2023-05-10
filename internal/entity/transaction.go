package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromAccountID *uint `json:"fromAccountID" gorm:"foreignkey:accountId"`
	ToAccountID   uint  `json:"toAccountID" gorm:"foreignkey:accountId;notNull"`

	Date        time.Time `json:"date" gorm:"notNull"`
	Amount      float64   `json:"amount" gorm:"notNull"`
	Description string    `json:"description"`
	Location    string    `json:"location" gorm:"size:128"`

	// `Transaction` belongs to `Category`, `CategoryID` is the foreign key
	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID"`
	// TODO: I suppose constrains here doesn't work - cascade
	Tags []Tag `gorm:"many2many:transaction_tags;constraint:OnDelete:CASCADE"`

	TransactionTypeID byte `json:"transactionTypeID" gorm:"notNull"`
	// TransactionType   TransactionType `gorm:"foreignKey:TransactionTypeID"`
}
