package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	UserID *uint
	Name   string `gorm:"notNull,size:64"`
	Color  string `gorm:"notNull;size:7"`
	Icon   string `gorm:"notNull;size:128"`
}
