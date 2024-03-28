package api

import (
	"fmt"
	"hotel-reservation/db"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {

	return func(c *fiber.Ctx) error {

		token, ok := c.GetReqHeaders()["X-Api-Token"]

		if !ok {
			return ErrorUnauthorized()
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrorUnauthorized()
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// fmt.Println("\n *** >>> invalid signing method", token.Header["alg"])
			return nil, ErrorUnauthorized()
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("\n*** >>> failed to parse token: ", err)
		return nil, ErrorUnauthorized()
	}

	if !token.Valid {
		fmt.Println("\n*** >>> [invalid token]")
		return nil, ErrorUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrorUnauthorized()
	}

	return claims, nil
}
