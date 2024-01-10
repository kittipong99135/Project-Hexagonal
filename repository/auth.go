package repository

import "auth-hex/models"

type AuthRepository interface {
	RepCreate(models.User) (*models.User, error)
	RepGetById(string) (*models.User, error)
	RepGetByEmail(string) (*models.User, int, error)
	RepUserExit(string) (int, error)
	RepSetRedis(models.Redis, int) error
	RepGetRedis(models.Redis) (*models.Redis, error)
	RepClearRedis(models.Redis) (*models.Redis, error)
}
