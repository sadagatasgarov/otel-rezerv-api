package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	//fmt.Println("isledi")
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	//fmt.Println(token)
	if !ok {
		return fmt.Errorf("1unauthorized")
	}
	claims, err := parseJWTToken(token[len(token)-1])
	if err != nil {
		return err
	}

	if claims["expires"]==nil {
		return fmt.Errorf("gecersiz token")
	}

	exFloat := claims["expires"].(float64)
	exp := int64(exFloat)
	if time.Now().Unix() > exp {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func parseJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method: ", token.Header["alg"])
			return nil, fmt.Errorf("unauhorizedd")
		}
		//secret:=os.Getenv("JWT_SECRET")
		secret := "salam"
		//fmt.Println("NEVERPRINTSECRET:", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse jwt:", err)
		return nil, fmt.Errorf("2unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("3unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		//fmt.Println(claims, "----------------")
		return nil, fmt.Errorf("4unauthorized")
	}

	return claims, nil
}
