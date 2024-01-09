package routes

import (
	"auth-hex/hadler"
	"auth-hex/middleware"
	"auth-hex/repository"
	"auth-hex/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB, rd *redis.Client) {
	api := app.Group("/api")

	// Auth Hexagonal
	authRep := repository.NewAuthRepository(db, rd)
	authSrv := service.NewAuthService(authRep)
	authHan := hadler.NewAuthHandler(authSrv)

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", authHan.Register)
	auth.Post("/login", authHan.Login)

	// User Hexagonal
	userRep := repository.NewUserRepository(db)
	userSrv := service.NewUserService(userRep)
	userHan := hadler.NewUserHandler(userSrv)

	// User Routes
	user := api.Group(
		"/user",
		middleware.RequestAuth(),
		middleware.RefreshAuth(),
	)
	user.Get("/params", userHan.UserParams)
	user.Get("/list", userHan.ListAllUser)
	user.Get("/read/:id", userHan.ReadUser)
	user.Put("/active/:id", userHan.Active)
	user.Put("/update/:id", userHan.Update)
	user.Delete("/remove/:id", userHan.Remove)

}
