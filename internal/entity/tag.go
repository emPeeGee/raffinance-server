package entity

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	UserID *uint
	Name   string `gorm:"notNull,size:64"`
	Color  string `gorm:"notNull;size:7"`
	Icon   string `gorm:"notNull;size:128"`
	// TODO: I suppose constrains here doesn't work - cascade
	Transactions []Transaction `gorm:"many2many:transaction_tags;constraint:OnDelete:RESTRICT"`
}
