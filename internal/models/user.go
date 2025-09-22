package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Email     string    `json:"email" gorm:"unique" validate:"required,email"`
	Name      string    `json:"name" validate:"required,min=1"`
	Age       int       `json:"age" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
