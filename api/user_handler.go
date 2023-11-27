package api

import (
	"context"
	"errors"

	db "hotel_api/storage"
	"hotel_api/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	return nil
}
func (h UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	user, err := h.userStore.DeleteUser(c.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"msg": "Silmek istediyiniz istifadeci tapilmadi"})
		}
		return err
	}
	msg := map[string]types.Users{}
	msg["silindi"] = *user
	return c.JSON(msg)
}

func (h UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return c.JSON(err)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"msg": "Axtardiginiz istifadeci tapilmadi"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
