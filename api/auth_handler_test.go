package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Jimmy",
		LastName:  "Jackme",
		Email:     "jimmy@mail.com",
		Password:  "fuqyourcouch",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthenticateWrongPassword(t *testing.T) {

	test_db := setup()
	defer test_db.teardown(t)

	insertedUser := insertTestUser(t, test_db.UserStore)
	insertedUser.EncryptedPassword = ""

	app := fiber.New()

	authHandler := NewAuthHandler(test_db.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jimmy@mail.com",
		Password: "fuqyourcouch-is-incorrect",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("\n*** >>> expected status 400 but got status: %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("\n*** >>> <genResp.type> should have been <error> - actual: <%s>", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("\n*** >>> <genResp.Msg> should have been <invalid credentials> - actual: <%s>", genResp.Type)
	}
}
func TestAuthenticateSuccess(t *testing.T) {

	test_db := setup()
	defer test_db.teardown(t)

	insertedUser := insertTestUser(t, test_db.UserStore)
	insertedUser.EncryptedPassword = ""

	app := fiber.New()

	authHandler := NewAuthHandler(test_db.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jimmy@mail.com",
		Password: "fuqyourcouch",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("\n*** >>> expected status 200 but got status: %d", resp.StatusCode)
	}

	var authResp AuthRespose
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	if (authResp.Token) == "" {
		t.Fatalf("\n*** >>> expected JWT token to be present in the auth response")
	}
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("\n*** >>> [insertedUser] mismatch with [authResp.user]")
	}
}
