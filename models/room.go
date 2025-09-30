package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;size:80;not null" json:"name"`
}
