package repository

import (
	"College/db"
	"College/models"
	"errors"

	"gorm.io/gorm"
)

func GetOrCreateRoom(name string) (models.Room, error) {
	var r models.Room
	err := db.DB.Where("name = ?", name).First(&r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r = models.Room{Name: name}
		err = db.DB.Create(&r).Error
	}
	return r, err
}

func FindRoomByName(name string) (models.Room, error) {
	var r models.Room
	err := db.DB.Where("name = ?", name).First(&r).Error
	return r, err
}

func ListRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := db.DB.Order("name asc").Find(&rooms).Error
	return rooms, err
}

// ListActiveRooms returns rooms that have at least one message, matching the
// Spring Boot side's DB-backed definition of "active". The previous
// implementation returned the in-memory hub's currently-connected rooms,
// which is a different, WS-session-scoped notion of "active" — a functional
// mismatch between the two apps that would have made /rooms/active results
// incomparable in K6 tests.
func ListActiveRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := db.DB.
		Joins("JOIN messages ON messages.room_id = rooms.id").
		Group("rooms.id").
		Order("rooms.name asc").
		Find(&rooms).Error
	return rooms, err
}
