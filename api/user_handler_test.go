package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

const testdburi = "mongodb://localhost:27017"
const testdbName = "hotel-reservation-test"

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testdb {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))

	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {

	test_db := setup()
	defer test_db.teardown(t)

	app := fiber.New()

	userHandler := NewUserHandler(test_db.UserStore)

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

	fmt.Println("\n*** >>> Response status: ", resp.Status)
}
