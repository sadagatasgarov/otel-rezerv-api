package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.Users)
	if !ok {
		return ErrUnAuthorized()
	}

	if !user.IsAdmin {
		return ErrUnAuthorized()
	}
	return c.Next()
}
