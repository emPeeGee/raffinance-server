package auth

import (
	"time"

	"github.com/lib/pq"
)

type createUserDTO struct {
	Password string `json:"password" validate:"required,min=4,max=256"`
	Email    string `json:"email" validate:"required,min=3,max=256"`
	Phone    string `json:"phone" validate:"required,min=1,max=16"`
	Username string `json:"username" validate:"required,min=3,max=64"`
	Name     string `json:"name" validate:"required,min=2,max=256"`
}

type credentialsDTO struct {
	Password string `json:"password" validate:"required,min=4,max=256"`
	Username string `json:"username" validate:"required,min=3,max=64"`
}

type userHashedPassword struct {
	Password string
}

type UserResponse struct {
	Username     string         `json:"username"`
	Name         string         `json:"name"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	LatestLogins pq.StringArray `json:"latestLogins" gorm:"type:varchar(64)[]"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
