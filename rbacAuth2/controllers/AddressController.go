package controllers

import "github.com/gofiber/fiber/v2"

func (uc *UserController) UpdateAddr(c *fiber.Ctx) error {
	return c.SendString("Admin Test Route")
}
