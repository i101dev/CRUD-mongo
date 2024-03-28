package api

import (
	"hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return ErrorUnauthorized()
	}

	if !user.IsAdmin {
		return ErrorUnauthorized()
	}

	return c.Next()
}
