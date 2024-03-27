package api

import (
	"errors"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"os"
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

type AuthRespose struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// -----------------------------------------------------------
// A handler should only do:
//   - serialization of the incoming request (JSON)
//   - database queries
//   - CALL business logic, rather than execute business logic internally
//   - return desired data
//
// .
func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {

	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}

		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthRespose{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {

	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))

	if err != nil {
		fmt.Println("\n*** >>> [error converting token to signed string] -", err)
	}

	return tokenStr
}
