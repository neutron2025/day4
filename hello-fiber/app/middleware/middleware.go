package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Logger() func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		fmt.Println("Request received")

		return c.Next()

	}

}
func SpecificUrlLogger() func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		fmt.Println("Request SpecificUrlLogger received")

		return c.Next()

	}

}
