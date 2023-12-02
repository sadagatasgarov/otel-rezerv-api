package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("isledi")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	fmt.Println("jjj: ",token)
	return nil
}
