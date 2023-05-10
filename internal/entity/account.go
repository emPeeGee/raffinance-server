package entity

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID   *uint
	Name     string `json:"name" gorm:"notNull;size:256"`
	Color    string `json:"color" gorm:"notNull;size:7"`
	Icon     string `gorm:"notNull;size:128"`
	Currency string `json:"currency" gorm:"notNull;size:10"`
}
