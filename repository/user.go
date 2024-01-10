package repository

import "auth-hex/models"

type UserRepository interface {
	RepGetAllUser() ([]models.User, error)
	RepGetUserById(string) (*models.User, error)
	RepActiveStatus(string, string) (*models.User, error)
	RepDeleteUser(string) (string, error)
	RepUpdateUser(string, *models.User) (string, error)
	RepUpdateExit(string, string) (int, error)
}
