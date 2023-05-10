package entity

type TransactionType struct {
	ID   byte   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"notNull"`
}
