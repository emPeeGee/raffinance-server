package category

import (
	"time"
)

type categoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CategoryShortResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type createCategoryDTO struct {
	Name  string `json:"name" validate:"required,min=2,max=64"`
	Color string `json:"color" validate:"required,hexcolor,min=7,max=7"`
	Icon  string `json:"icon" validate:"required,max=128"`
}

type updateCategoryDTO struct {
	Name  string `json:"name" validate:"required,min=2,max=64"`
	Color string `json:"color" validate:"required,hexcolor,min=7,max=7"`
	Icon  string `json:"icon" validate:"required,max=128"`
}
