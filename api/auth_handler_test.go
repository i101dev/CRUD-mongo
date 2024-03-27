package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-reservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateWrongPassword(t *testing.T) {

	test_db := setup(t)
	defer test_db.teardown(t)

	fixtures.AddUser(test_db.Store, "test", "user", false)

	app := fiber.New()

	authHandler := NewAuthHandler(test_db.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "test@user.com",
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

	test_db := setup(t)
	defer test_db.teardown(t)

	// insertedUser := insertTestUser(t, test_db.User)
	insertedUser := fixtures.AddUser(test_db.Store, "test", "user", false)
	insertedUser.EncryptedPassword = ""

	app := fiber.New()

	authHandler := NewAuthHandler(test_db.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    fmt.Sprintf("%s@%s.com", insertedUser.FirstName, insertedUser.LastName),
		Password: fmt.Sprintf("%s_%s", insertedUser.FirstName, insertedUser.LastName),
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
