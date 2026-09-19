package services

import (
	"College/models"
	"College/repository"
)

func EnsureRoom(name string) (models.Room, error) {
	return repository.GetOrCreateRoom(name)
}

func FindRoomByName(name string) (models.Room, error) {
	return repository.FindRoomByName(name)
}

func ListRooms() ([]models.Room, error) {
	return repository.ListRooms()
}

func ListActiveRooms() ([]models.Room, error) {
	return repository.ListActiveRooms()
}
