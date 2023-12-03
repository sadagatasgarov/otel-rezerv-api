package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("isledi")

	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if err := parseJWTToken(token[len(token)-1]); err != nil {
		return err
	}

	return nil
}

func parseJWTToken(tokenStr string) error {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method: ", token.Header["alg"])
			return nil, fmt.Errorf("unauhorizedd")
		}
		//secret:=os.Getenv("JWT_SECRET")
		secret := "Bugizlisozdur"
		fmt.Println("NEVERPRINTSECRET:", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse jwt:", err)

		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["foo"], claims["nbf"])
		fmt.Println(claims)
	}
	return fmt.Errorf("unauhorizedd")
}
