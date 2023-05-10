package entity

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID *uint
	Name   string `json:"name" gorm:"notNull;size:256"`
	Phone  string `json:"phone" gorm:"size:16"`
	Email  string `json:"email" gorm:"notNull;size:256"`
}
