package middleware

import (
	"blog-auth-server/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var SecretKey = []byte("SecretKey")

type Middleware struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewMiddleware(ctx context.Context, redisClient *redis.Client) *Middleware {
	return &Middleware{
		ctx:         ctx,
		redisClient: redisClient,
	}
}
func (uc *Middleware) AdminMiddlewareHandler(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	entry := c.Get("Entry")
	entryInt, err := strconv.Atoi(entry)
	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		log.Println("Token parsing error:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid token"})
	}
	hash := ""
	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			hash = claims["hash"].(string)
			expirationFloat := claims["exp"].(float64)
			if !ok {
				log.Println("Expiration not found in claims")
				return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
			}
			expiration := time.Unix(int64(expirationFloat), 0)

			// Compare the expiration time with the current time
			if time.Now().After(expiration) {
				log.Println("Token has expired")
				return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
			}
		} else {
			log.Println("Invalid claims")
			return err
		}
	} else {
		log.Println("Token is invalid")
	}
	value, err := uc.redisClient.Get(uc.ctx, hash).Result()
	var retrievedPermissions []models.Permissions
	if err := json.Unmarshal([]byte(value), &retrievedPermissions); err != nil {
		log.Println("Error deserializing permissions:", err)
		return err
	}
	fmt.Println(retrievedPermissions)
	for _, permission := range retrievedPermissions {
		if permission.Entry == entryInt {
			if permission.AdminFlag {
				return c.Next()
			} else {
				return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
			}
		} else {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
	}
	return nil
}

func (uc *Middleware) UserMiddlewareHandler(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if authorization == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}
	entry := c.Get("Entry")
	entryInt, err := strconv.Atoi(entry)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Entry format"})
	}

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		log.Println("Token parsing error:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid token"})
	}
	hash := ""
	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			hash = claims["hash"].(string)
			expirationFloat := claims["exp"].(float64)
			if !ok {
				log.Println("Expiration not found in claims")
				return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
			}
			expiration := time.Unix(int64(expirationFloat), 0)

			// Compare the expiration time with the current time
			if time.Now().After(expiration) {
				log.Println("Token has expired")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
			}
		} else {
			log.Println("Invalid claims")
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid claims"})
		}
	} else {
		log.Println("Token is invalid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is invalid"})
	}
	value, err := uc.redisClient.Get(uc.ctx, hash).Result()
	if err != nil {
		log.Println("Redis retrieval error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving permissions"})
	}
	var retrievedPermissions []models.Permissions
	if err := json.Unmarshal([]byte(value), &retrievedPermissions); err != nil {
		log.Println("Error deserializing permissions:", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deserializing permissions"})
	}
	fmt.Println(retrievedPermissions)

	for _, permission := range retrievedPermissions {
		if permission.Entry == entryInt {
			return c.Next()
		} else {
			log.Println("Permission denied for Entry:", entryInt)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}
	}
	return nil
}
