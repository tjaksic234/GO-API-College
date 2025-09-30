package repository

import (
	"College/db"
	"College/models"
)

func SaveMessage(m *models.Message) error {
	return db.DB.Create(m).Error
}

func GetMessagesByRoomID(roomID uint, limit, offset int) ([]models.Message, error) {
	var msgs []models.Message
	q := db.DB.
		Preload("User").
		Preload("Room").
		Where("room_id = ?", roomID).
		Order("created_at asc")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	err := q.Find(&msgs).Error
	return msgs, err
}
