package repository

import (
	"College/db"
	"College/models"
	"errors"

	"gorm.io/gorm"
)

func GetOrCreateUser(username string) (models.User, error) {
	var u models.User
	err := db.DB.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u = models.User{Username: username}
		err = db.DB.Create(&u).Error
	}
	return u, err
}

func ListUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Order("username asc").Find(&users).Error
	return users, err
}
