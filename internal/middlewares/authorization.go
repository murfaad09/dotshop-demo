package middlewares

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/config"
)

func NewAuthorizationMiddleware() fiber.Handler {
	c := config.GetConfig()

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(c.JWTSecret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		},
	})
}
