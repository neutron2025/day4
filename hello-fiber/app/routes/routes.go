package routes

import (
	"hello-fiber/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})
	app.Get("/protected", middleware.SpecificUrlLogger(), func(c *fiber.Ctx) error {
		return c.SendString("This route is protected by middleware.SpecificUrlLogger()")
	})
	// elegant way to handle request chain errors using middleware
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Handle the error here
				c.Status(fiber.StatusInternalServerError).SendString("Something went wrong!")
			}
		}()
		return c.Next()
	})

}
