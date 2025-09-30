package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	RoomID  uint   `json:"room_id"`
	UserID  uint   `json:"user_id"`
	Content string `gorm:"type:text" json:"content"`

	Room Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
