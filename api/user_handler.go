package api

import (
	"errors"
	"hotel-reservation/db"
	"hotel-reservation/types"

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

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {

	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := db.Map{"_id": userID}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {

	userID := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted": userID})
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {

	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	intsertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(intsertedUser)
}
func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {

	id := c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "document not found"})
		}

		return err
	}

	return c.JSON(user)
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	users, err := h.userStore.GetUsers(c.Context())

	if err != nil {
		return ErrorResourceNotFound("user data")
	}

	return c.JSON(users)
}
