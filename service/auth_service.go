package service

import (
	"auth-hex/models"
	"auth-hex/repository"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	act_token, err := CreateToken(result, "JWT_SECRET", 30)
	if err != nil {
		return nil, errors.New("Create accesstoke invalids")
	}
	setRedis := models.Redis{
		Key:   "access_token:" + strconv.Itoa(int(result.ID)),
		Value: act_token,
	}
	s.authRep.RepSetRedis(setRedis, 0)
	accessRedis, _ := s.authRep.RepGetRedis(setRedis)

	rfh_token, err := CreateRefreshToken(result, "JWT_REFRESH")

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

func (s authService) SrvValidate(tokenStr models.LoginResponse) (*models.LoginResponse, error) {
	access_token := strings.TrimPrefix(tokenStr.Access_token, "Bearer ")
	result, err := jwt.Parse(access_token, TestValids)

	var uid string
	var token_expired bool

	if !result.Valid {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			uid, err = result.Claims.GetIssuer()
			token_expired = true
		default:
			token_expired = false
			return nil, errors.New("error handle this token:")
		}
	}

	setRedis := models.Redis{
		Key: "refresh_token:" + uid,
	}

	refreshRedis, _ := s.authRep.RepGetRedis(setRedis)
	if token_expired {
		refresh_valid, err := jwt.Parse(refreshRedis.Value, TestValids)
		fmt.Println(refresh_valid)

		if !refresh_valid.Valid {
			return nil, errors.New("refresh token invalid")
		}

		userResultValid, err := s.authRep.RepGetById(uid)
		if err != nil {
			return nil, errors.New("get user invalid")
		}

		act_token, err := CreateToken(userResultValid, "JWT_SECRET", 30)
		setRedis := models.Redis{
			Key:   "access_token:" + uid,
			Value: act_token,
		}

		s.authRep.RepSetRedis(setRedis, 0)
		accessRedis, _ := s.authRep.RepGetRedis(setRedis)
		_ = accessRedis

		resultLogin := models.LoginResponse{
			Status:        "refresh",
			Access_token:  act_token,
			Refresh_token: refreshRedis.Value,
		}
		return &resultLogin, errors.New("refresh token")
	}

	resultLogin := models.LoginResponse{
		Status:        "success",
		Access_token:  access_token,
		Refresh_token: refreshRedis.Value,
	}
	return &resultLogin, nil
}

func CreateToken(userResult *models.User, env string, exp int) (string, error) {
	cliams := jwt.MapClaims{
		"iss":    strconv.Itoa(int(userResult.ID)),
		"id":     userResult.ID,
		"name":   userResult.Name,
		"email":  userResult.Email,
		"role":   userResult.Role,
		"status": userResult.Status,
		"rank":   userResult.Rank,
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}

func CreateRefreshToken(userResult *models.User, env string) (string, error) {
	cliams := jwt.MapClaims{
		"iss": strconv.Itoa(int(userResult.ID)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}

func TestValids(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}
