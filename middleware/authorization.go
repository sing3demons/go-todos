package middleware

import (
	"github.com/casbin/casbin"
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

func Authorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := c.Locals("sub")
		if sub == nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		path := c.Request().URI().Path()

		method := c.Request().Header.Method()

		ok := enforcer.Enforce(sub.(model.User), string(path), string(method))

		if !ok {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.Next()
	}
}
