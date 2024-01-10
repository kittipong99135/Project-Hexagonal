package service

import (
	"auth-hex/models"

	_ "github.com/golang-jwt/jwt"
)

type AuthService interface {
	SrvRegister(models.UserRequest) (*models.UserResponse, error)
	SrvLogin(models.AuthRequest) (*models.LoginResponse, error)
	SrvValidate(models.LoginResponse) (*models.LoginResponse, error)
}
