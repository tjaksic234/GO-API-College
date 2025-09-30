package services

import (
	"College/models"
	"College/repository"
)

func SaveMessage(roomID uint, userID uint, content string) error {
	msg := models.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	}
	return repository.SaveMessage(&msg)
}

func GetRoomMessages(roomID uint, limit, offset int) ([]models.Message, error) {
	return repository.GetMessagesByRoomID(roomID, limit, offset)
}
