package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartController struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCartController(collection *mongo.Collection, ctx context.Context) *CartController {
	return &CartController{
		collection: collection,
		ctx:        ctx,
	}
}

func (uc *CartController) AddtoCart(c *fiber.Ctx) error {
	return c.SendString(" Test Route")
}

func (uc *CartController) DelfromCart(c *fiber.Ctx) error {
	return c.SendString(" Test Route")
}
func (uc *CartController) AllfromCart(c *fiber.Ctx) error {
	return c.SendString(" Test Route")
}

func (uc *CartController) UpdatefromCart(c *fiber.Ctx) error {
	return c.SendString(" Test Route")
}
