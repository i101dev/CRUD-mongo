package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {

	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("not authorized")
	}

	claims, err := validateToken(token[0])
	if err != nil {
		return err
	}

	// expiresStr, ok := claims["expires"].(string)
	// if !ok {
	// 	return fmt.Errorf("expires claim is not a string")
	// }

	// expires, err := time.Parse(time.RFC3339, expiresStr)
	// if err != nil {
	// 	return fmt.Errorf("error parsing expiration time: %v", err)
	// }

	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)

	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("\n *** >>> invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		// fmt.Println("\n*** >>> [secret] - ", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("\n*** >>> failed to parse token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("\n*** >>> [invalid token]")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
