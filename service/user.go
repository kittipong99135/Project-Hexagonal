package service

import (
	"auth-hex/models"
)

type UserService interface {
	SrvGetAllUser() ([]models.UserResponse, error)
	SrvGetUserById(string) (*models.UserResponse, error)
	SrvActiveUser(string) (*models.UserResponse, error)
	SrvDeleteUser(string) (string, error)
	SrvUpdateUser(string, *models.User) (string, error)
}
