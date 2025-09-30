package services

import (
	"College/models"
	"College/repository"
)

func EnsureUser(username string) (models.User, error) {
	return repository.GetOrCreateUser(username)
}

func ListUsers() ([]models.User, error) {
	return repository.ListUsers()
}
