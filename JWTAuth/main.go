package main

import (
	"jwt-auth/config"
	"jwt-auth/database"
	"jwt-auth/handlers"
	"jwt-auth/middleware"
	"jwt-auth/models"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New()

	database.Connect(cfg)
	models.Migrate(database.DB)

	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	app.Use(middleware.Auth)
	app.Get("/user", handlers.User)
	app.Post("/logout", handlers.Logout)
	app.Post("/refresh", handlers.Refresh)

	app.Listen(":6000")
}
