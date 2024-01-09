package service

import (
	"auth-hex/models"
	"auth-hex/repository"
	"errors"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRep repository.AuthRepository
}

func NewAuthService(authRep repository.AuthRepository) AuthService {
	return authService{authRep: authRep}
}

func (s authService) SrvRegister(user models.UserRequest) (*models.UserResponse, error) {
	rows, err := s.authRep.RepUserExit(user.Email)
	if err != nil {
		return nil, errors.New("Register invalids!")
	}
	if rows != 0 {
		return nil, errors.New("Email exists!")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, errors.New("Invalid hash passward")
	}

	userRegis := models.User{
		Email:    user.Email,
		Password: string(hash),
		Name:     user.Name,
		Age:      user.Age,
		Phone:    user.Phone,
		Rank:     user.Rank,
		Role:     "user",
		Status:   "nactive",
	}

	result, err := s.authRep.RepCreate(userRegis)
	if err != nil {
		return nil, errors.New("Invalid register!")
	}

	userResponse := models.UserResponse{
		Email:  result.Email,
		Name:   result.Name,
		Phone:  result.Phone,
		Rank:   result.Rank,
		Status: result.Status,
		Role:   result.Role,
	}

	return &userResponse, nil
}

func (s authService) SrvLogin(login models.AuthRequest) (*models.LoginResponse, error) {
	result, rows, err := s.authRep.RepGetByEmail(login.Email)
	if err != nil {
		return nil, errors.New("Register invalids!")
	}
	if rows == 0 {
		return nil, errors.New("Email invalid!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(login.Password))
	if err != nil {
		return nil, errors.New("Password invalid")
	}

	act_token, err := CreateToken(result, "JWT_SECRET")
	if err != nil {
		return nil, errors.New("Create accesstoke invalids")
	}
	setRedis := models.Redis{
		Key:   "access_token:" + strconv.Itoa(int(result.ID)),
		Value: act_token,
	}
	s.authRep.RepSetRedis(setRedis, 2)
	accessRedis, _ := s.authRep.RepGetRedis(setRedis)

	rfh_token, err := CreateToken(result, "JWT_SECRET")
	if err != nil {
		return nil, errors.New("Create accesstoke invalids")
	}
	setRedis = models.Redis{
		Key:   "refresh_token:" + strconv.Itoa(int(result.ID)),
		Value: rfh_token,
	}
	s.authRep.RepSetRedis(setRedis, 0)
	refreshRedis, _ := s.authRep.RepGetRedis(setRedis)

	resultLogin := models.LoginResponse{
		Status:        "success",
		Access_token:  accessRedis.Value,
		Refresh_token: refreshRedis.Value,
	}

	return &resultLogin, nil
}

func CreateToken(userResult *models.User, env string) (string, error) {
	cliams := jwt.MapClaims{
		"uid":    userResult.ID,
		"name":   userResult.Name,
		"email":  userResult.Email,
		"role":   userResult.Role,
		"status": userResult.Status,
		"rank":   userResult.Rank,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}
