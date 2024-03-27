package api

import (
	"bytes"
	"encoding/json"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {

	test_db := setup(t)
	defer test_db.teardown(t)

	app := fiber.New()

	userHandler := NewUserHandler(test_db.User)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "jimmy@gmail.com",
		FirstName: "Jimmy",
		LastName:  "Tru",
		Password:  "123qweasdzxc",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	var user types.User

	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting a user ID to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the encrypted password to not be included in the response")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("FirstName: expected %s, but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("LastName: expected %s, but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("Email: expected %s, but got %s", params.Email, user.Email)
	}

}
