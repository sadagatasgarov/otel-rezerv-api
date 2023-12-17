package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"

	"github.com/gofiber/fiber/v2"

	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("token not present in the header")
			return ErrUnAuthorized()
		}

		claims, err := parseJWTToken(token[len(token)-1])
		if err != nil {
			return err
		}

		if claims["expires"] == nil {
			return fmt.Errorf("token gecersiz")
		}

		exFloat := claims["expires"].(float64)
		exp := int64(exFloat)
		if time.Now().Unix() > exp {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func parseJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method: ", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		//secret := "salam"
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse jwt:", err)
		return nil, ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return nil, ErrUnAuthorized()
	}

	return claims, nil
}
