package middleware

import (
	"auth-hex/database"
	"context"
	"fmt"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsRefreshToken struct {
	Uid    string
	Name   string
	Email  string
	Role   string
	Status string
	Rank   string
}

func RequestAuth() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		// ErrorHandler: errNext,
	})
}

func RefreshAuth() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_REFRESH"))},
		SuccessHandler: resendToken,
		ErrorHandler:   errNext,
	})
}

func errNext(c *fiber.Ctx, err error) error {
	c.Next()
	return nil
}

func resendToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])
	name := fmt.Sprintf("%v", claims["name"])
	email := fmt.Sprintf("%v", claims["email"])
	role := fmt.Sprintf("%v", claims["role"])
	status := fmt.Sprintf("%v", claims["status"])
	rank := fmt.Sprintf("%v", claims["rank"])
	getAccessToken := GetToken("access_token:" + uid)

	if getAccessToken != "" {
		fmt.Println("have token: \n" + getAccessToken + "\n")
		c.Next()
		return nil
	}

	claimsRefreshToken := ClaimsRefreshToken{
		Uid:    uid,
		Name:   name,
		Email:  email,
		Role:   role,
		Status: status,
		Rank:   rank,
	}

	acc_token, err := RefreshToken(claimsRefreshToken, "JWT_SECRET")
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Refresh token invalid.",
			"error":   err.Error(),
		})
	}
	fmt.Println("refresh token: \n" + acc_token + "\n:")
	SetAccessToken("access_token:"+uid, acc_token)

	c.Next()
	return c.Status(200).JSON(fiber.Map{
		"status":       "warning",
		"message":      "Warning : Refresh token.",
		"Access_token": GetToken("access_token:" + uid),
	})
}

func RefreshToken(claimsRefreshToken ClaimsRefreshToken, env string) (string, error) {
	cliams := jwt.MapClaims{
		"uid":    claimsRefreshToken.Uid,
		"name":   claimsRefreshToken.Name,
		"email":  claimsRefreshToken.Email,
		"role":   claimsRefreshToken.Role,
		"status": claimsRefreshToken.Status,
		"rank":   claimsRefreshToken.Rank,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}

func SetAccessToken(key string, token string) {
	rd := database.RD_Init()
	ctx := context.Background()
	rd.Set(ctx, key, token, time.Second*10)
}

func GetToken(key string) string {
	rd := database.RD_Init()
	ctx := context.Background()
	val, _ := rd.Get(ctx, key).Result()
	return val
}
