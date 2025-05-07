package utils

import "github.com/gofiber/fiber/v2"

func GetUserID(c *fiber.Ctx) (uint, error) {
	val := c.Locals("user_id")
	userID, ok := val.(uint)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return userID, nil
}
