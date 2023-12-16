package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"
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
