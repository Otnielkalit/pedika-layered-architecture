package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractToken(c *fiber.Ctx) (*jwt.Token, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token tidak ditemukan")
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Format token tidak valid")
	}
	tokenString := splitToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token tidak valid")
	}
	return token, nil
}

func AdminMiddleware(c *fiber.Ctx) error {
	token, err := ExtractToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses hanya untuk admin",
			"data":    nil,
		})
	}
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	return c.Next()
}

func MasyarakatMiddleware(c *fiber.Ctx) error {
	token, err := ExtractToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if claims["role"] != "masyarakat" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses hanya untuk masyarakat",
			"data":    nil,
		})
	}
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	return c.Next()
}
