package middleware

import (
	"fmt"
	"gitlab.com/sadagatasgarov/otel-rezervasiya-api/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.Users)
	if !ok {
		return fmt.Errorf("not authorized")
	}

	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}
	return c.Next()
}
