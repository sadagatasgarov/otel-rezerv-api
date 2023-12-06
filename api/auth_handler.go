package api

import (
	"errors"
	"fmt"
	"net/http"
	db "sadagatasgarov/hotel_rezerv_api/storage"
	"sadagatasgarov/hotel_rezerv_api/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.Users `json:"user"`
	Token string       `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msq"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	//fmt.Println(params)

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "Istifadeci yoxdut tapilmadi"})
		}
		return invalidCredentials(c)
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	//fmt.Println("authenticated->", user)
	return c.JSON(resp)
}

func createTokenFromUser(user *types.Users) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"validTill": validTill,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("salam")

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Token with secret erroru", err)
	}
	return tokenStr
}
