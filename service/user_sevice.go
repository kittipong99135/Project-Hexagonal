package service

import (
	"auth-hex/models"
	"auth-hex/repository"
	"errors"
	"fmt"
	"strconv"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) SrvGetAllUser() ([]models.UserResponse, error) {
	users, err := s.userRepo.RepGetAllUser()
	if err != nil {
		return nil, errors.New("Get user invelid.")
	}

	userResponses := []models.UserResponse{}

	for _, user := range users {
		userRespose := models.UserResponse{
			ID:     strconv.Itoa(int(user.ID)),
			Email:  user.Email,
			Name:   user.Name,
			Phone:  user.Phone,
			Rank:   user.Rank,
			Status: user.Status,
			Role:   user.Role,
		}
		userResponses = append(userResponses, userRespose)
	}
	return userResponses, nil
}

func (s userService) SrvGetUserById(uid string) (*models.UserResponse, error) {
	user, err := s.userRepo.RepGetUserById(uid)
	if err != nil {
		return nil, errors.New("Get user invelid.")
	}
	userResponse := models.UserResponse{
		ID:     strconv.Itoa(int(user.ID)),
		Email:  user.Email,
		Name:   user.Name,
		Phone:  user.Phone,
		Rank:   user.Rank,
		Status: user.Status,
		Role:   user.Role,
	}

	return &userResponse, nil
}

func (s userService) SrvActiveUser(uid string) (*models.UserResponse, error) {
	user, err := s.userRepo.RepGetUserById(uid)
	if err != nil {
		return nil, errors.New("Get user invalid.")
	}

	resultActive, err := s.userRepo.RepActiveStatus(uid, user.Status)
	if err != nil {
		return nil, errors.New("Active user invalid.")
	}
	responseActive := models.UserResponse{
		ID:     strconv.Itoa(int(resultActive.ID)),
		Email:  resultActive.Email,
		Name:   resultActive.Name,
		Phone:  resultActive.Phone,
		Rank:   resultActive.Rank,
		Status: resultActive.Status,
		Role:   resultActive.Role,
	}

	return &responseActive, nil
}

func (s userService) SrvDeleteUser(uid string) (string, error) {
	result, err := s.userRepo.RepDeleteUser(uid)
	if err != nil {
		return "", errors.New("Delete user invalid.")
	}

	return result, nil
}

func (s userService) SrvUpdateUser(uid string, userUpdate *models.User) (string, error) {

	emailExit, err := s.userRepo.RepUpdateExit(uid, userUpdate.Email)
	if err != nil {
		fmt.Println("Invalids bro")
		return "", errors.New("Update user invalid.")
	}
	if emailExit != 0 {
		return "", errors.New("Email exists.")
	}

	result, err := s.userRepo.RepUpdateUser(uid, userUpdate)

	if err != nil {
		return "", errors.New("Update user invalid.")
	}
	return result, nil
}
