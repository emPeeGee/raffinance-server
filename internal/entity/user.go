package entity

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string         `json:"name" gorm:"notNull;size:256"`
	Email        string         `json:"email" gorm:"notNull;unique;size:256"`
	Username     string         `json:"username" gorm:"notNull;unique;size:64"`
	Password     string         `json:"password" gorm:"notNull;size:256"`
	Phone        string         `json:"phone" gorm:"size:16"`
	LatestLogins pq.StringArray `json:"latestLogins" gorm:"type:varchar(64)[]"`
	Accounts     []Account
	Contacts     []Contact
	Categories   []Category
	Tags         []Tag
}
