package main

import (
	"auth-hex/database"
	"auth-hex/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	db := database.DB_Init()
	rd := database.RD_Init()

	app := fiber.New()
	// Initialize default config
	app.Use(cors.New())
	routes.Routes(app, db, rd)

	app.Listen("127.0.0.1:3000")
}
