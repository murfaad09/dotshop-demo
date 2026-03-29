package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/harishash/dotshop-be/internal/config"
)

var jwtKey = []byte(config.GetConfig().JWTSecret)

type Claims struct {
	Email  string `json:"email"`
	RoleId uint64 `json:"role_id"`
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string, roleID uint64, userId uint64) (string, error) {
	var expirationTime time.Time
	if roleID == 1 {
		expirationTime = time.Now().Add(20 * 24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(24 * time.Hour)
	}
	claims := &Claims{
		Email:  email,
		RoleId: roleID,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleId)
		c.Locals("user_id", claims.UserId)

		return c.Next()
	}
}

func JWTMiddlewareValidation(token string) bool {
	if token == "" {
		return false
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return false
	}

	return true
}

func GenerateJWTForForgotPassword(email string, roleID uint64, userId uint64) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Email:  email,
		RoleId: roleID,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
