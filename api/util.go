package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezervasiya-api/types"
)

func getAuthUser(c *fiber.Ctx) (*types.Users, error) {
	user, ok := c.Context().UserValue("user").(*types.Users)
	if !ok {
		return nil, fmt.Errorf("user ok deyil")
	}
	return user, nil
}
