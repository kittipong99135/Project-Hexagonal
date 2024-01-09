package hadler

import (
	"auth-hex/models"
	"auth-hex/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) UserParams(c *fiber.Ctx) error {
	type paramsUser struct {
		Uid    string
		Name   string
		Email  string
		Role   string
		Status string
		Rank   string
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])
	name := fmt.Sprintf("%v", claims["name"])
	email := fmt.Sprintf("%v", claims["email"])
	role := fmt.Sprintf("%v", claims["role"])
	status := fmt.Sprintf("%v", claims["status"])
	rank := fmt.Sprintf("%v", claims["rank"])
	params := paramsUser{
		Uid:    uid,
		Name:   name,
		Email:  email,
		Role:   role,
		Status: status,
		Rank:   rank,
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Set params user success.",
		"params":  params,
	})
}

func (h userHandler) ListAllUser(c *fiber.Ctx) error {
	result, err := h.userSrv.SrvGetAllUser()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "List user invalids",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "List user success",
		"result":  result,
	})
}

func (h userHandler) ReadUser(c *fiber.Ctx) error {
	uid := c.Params("id")
	result, err := h.userSrv.SrvGetUserById(uid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Read user invalids",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "List user success",
		"result":  result,
	})
}

func (h userHandler) Active(c *fiber.Ctx) error {
	uid := c.Params("id")
	result, err := h.userSrv.SrvActiveUser(uid)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Read user invalids",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "List user success",
		"result":  result,
	})
}

func (h userHandler) Remove(c *fiber.Ctx) error {
	uid := c.Params("id")
	result, err := h.userSrv.SrvDeleteUser(uid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Delete user invalids",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Delete user success",
		"result":  result,
	})
}

func (h userHandler) Update(c *fiber.Ctx) error {
	uid := c.Params("id")
	user := models.User{}
	c.BodyParser(&user)
	result, err := h.userSrv.SrvUpdateUser(uid, &user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Update user invalids",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Update user success",
		"result":  result,
	})
}
