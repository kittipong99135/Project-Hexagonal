package repository

import (
	"auth-hex/models"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
	rd *redis.Client
}

func NewAuthRepository(db *gorm.DB, rd *redis.Client) AuthRepository {
	return authRepository{db: db, rd: rd}
}

func (r authRepository) RepCreate(user models.User) (*models.User, error) {
	result := r.db.Create(&user)
	return &user, result.Error
}

func (r authRepository) RepGetById(id int) (*models.User, error) {
	var user models.User
	result := r.db.Find(&user, "id =?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, result.Error
}

func (r authRepository) RepGetByEmail(email string) (*models.User, int, error) {
	var user models.User
	result := r.db.Find(&user, "email =?", email)
	if result.RowsAffected != 1 {
		return nil, int(result.RowsAffected), result.Error
	}
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return &user, int(result.RowsAffected), result.Error
}

func (r authRepository) RepUserExit(email string) (int, error) {
	var userExist models.User
	result := r.db.Find(&userExist, "email =?", email)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r authRepository) RepSetRedis(redis models.Redis, exp int) error {
	ctx := context.Background()
	result := r.rd.Set(ctx, redis.Key, redis.Value, time.Hour*time.Duration(exp))
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (r authRepository) RepGetRedis(redis models.Redis) (*models.Redis, error) {
	ctx := context.Background()
	val, _ := r.rd.Get(ctx, redis.Key).Result()
	result := models.Redis{
		Key:   redis.Key,
		Value: val,
	}
	return &result, nil
}

func (r authRepository) RepClearRedis(models.Redis) (*models.Redis, error) {
	return nil, nil
}
