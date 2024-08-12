package main

import (
	"hello-fiber/app/middleware"
	"hello-fiber/app/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Use(middleware.Logger()) //applying middleware globally
	routes.SetupRoutes(app)

	// databaseURL := config.GetDatabaseURL()
	app.Listen(":3000")
}
