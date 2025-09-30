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
