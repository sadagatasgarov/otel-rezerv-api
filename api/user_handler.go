package api

import (
	"context"
	db "hotel_api/storage"
	"hotel_api/types"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.Users{
		FirstName: "sada",
		LastName:  "asga",
	}
	return c.JSON(u)
}

func (h UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
