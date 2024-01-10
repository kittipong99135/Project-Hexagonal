package hadler

import (
	"auth-hex/models"
	"auth-hex/service"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authSrv service.AuthService
}

func NewAuthHandler(authSrv service.AuthService) authHandler {
	return authHandler{authSrv: authSrv}
}

func (h authHandler) Register(c *fiber.Ctx) error {

	regisReq := models.UserRequest{}
	err := c.BodyParser(&regisReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "Invalid",
		})
	}
	resultReq, err := h.authSrv.SrvRegister(regisReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "Invalid",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success : Register user success.",
		"user":    resultReq,
	})
}

func (h authHandler) Login(c *fiber.Ctx) error {

	loginReq := models.AuthRequest{}
	err := c.BodyParser(&loginReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
		})
	}
	resultReq, err := h.authSrv.SrvLogin(loginReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
		})
	}

	return c.Status(200).JSON(resultReq)
}

func (h authHandler) ValidatetEst(c *fiber.Ctx) error {
	access_token := c.Get("Authorization")

	set_validate := models.LoginResponse{
		Access_token: access_token,
	}
	resultToken, err := h.authSrv.SrvValidate(set_validate)

	if err != nil {
		if err.Error() == "refresh token" {
			return c.Status(200).JSON(fiber.Map{
				"status":  "warning",
				"message": "Warning : Token has expired",
				"result":  resultToken,
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Token not validate",
			"error":   err.Error(),
		})
	}

	c.Next()
	return nil
}
