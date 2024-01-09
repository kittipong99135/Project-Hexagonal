package main

import (
	"auth-hex/database"
	"auth-hex/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db := database.DB_Init()
	rd := database.RD_Init()

	app := fiber.New()
	routes.Routes(app, db, rd)

	app.Listen("127.0.0.1:3000")
}
