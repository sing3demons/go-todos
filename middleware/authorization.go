package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/model"
)

func Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := c.Locals("sub").(model.User)
		if sub.Role == "Admin" || sub.Role == "Editor" {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusForbidden)
	}
}
