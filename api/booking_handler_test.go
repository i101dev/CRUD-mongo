package api

import (
	"encoding/json"
	"fmt"
	"hotel-reservation/db/fixtures"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {

	db := setup(t)

	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "franky", "twoshoes", false)
		user           = fixtures.AddUser(db.Store, "jimmy", "jackme", false)
		hotel          = fixtures.AddHotel(db.Store, "Hotel Bravo", "Semple", 88, nil)
		room           = fixtures.AddRoom(db.Store, "medium", true, 430, hotel.ID)
		from           = time.Now().AddDate(0, 0, 1)
		till           = from.AddDate(0, 0, 3)
		numPersons     = 2
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till, numPersons)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		route          = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 status code: %d", resp.StatusCode)
	}

	var bookingResp *types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("\n*** >>> [booking.ID] mismatch - expected %s; received: %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("\n*** >>> [booking.UserID] mismatch - expected %s; received: %s", booking.UserID, bookingResp.UserID)
	}

	// Test non-booking-user access -------------------------------------------------
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("\n*** >>> expected non-admin request to fail (not 200). status: %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {

	db := setup(t)

	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "hotel", "admin", true)
		user           = fixtures.AddUser(db.Store, "jimmy", "jackme", false)
		hotel          = fixtures.AddHotel(db.Store, "Hotel Bravo", "Semple", 88, nil)
		room           = fixtures.AddRoom(db.Store, "medium", true, 430, hotel.ID)
		from           = time.Now().AddDate(0, 0, 1)
		till           = from.AddDate(0, 0, 3)
		numPersons     = 2
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till, numPersons)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response - %d", resp.StatusCode)
	}

	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking - received: (%d) bookings", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("\n*** >>> [booking.ID] mismatch - expected %s; received: %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("\n*** >>> [booking.UserID] mismatch - expected %s; received: %s", booking.UserID, have.UserID)
	}

	// Test non-admin user access -------------------------------------------------
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("\n*** >>> expected [http.StatusUnauthorized] - received: %d", resp.StatusCode)
	}
}
