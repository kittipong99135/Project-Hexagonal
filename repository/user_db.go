package repository

import (
	"auth-hex/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{db: db}
}

func (r userRepository) RepGetAllUser() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r userRepository) RepGetUserById(uid string) (*models.User, error) {
	var user models.User
	result := r.db.Find(&user, "id =?", uid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, result.Error
}

func (r userRepository) RepActiveStatus(uid string, status string) (*models.User, error) {
	var activeUser models.User
	if status != "active" {
		activeUser = models.User{Status: "active"}
	} else {
		activeUser = models.User{Status: "nactive"}
	}

	r.db.Where("id = ?", uid).Updates(&activeUser)

	return &activeUser, nil
}

func (r userRepository) RepDeleteUser(uid string) (string, error) {
	user := models.User{}
	result := r.db.Delete(&user, "id = ?", uid)
	if result.Error != nil {
		return "", result.Error
	}
	return "Delete user success.", nil
}

func (r userRepository) RepUpdateUser(uid string, userUpdate *models.User) (string, error) {
	result := r.db.Where("id = ?", uid).Updates(&userUpdate)
	if result.Error != nil {
		return "", result.Error
	}
	return "Update user success.", nil
}
