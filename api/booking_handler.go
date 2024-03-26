package api

import (
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(s *db.Store) *BookingHandler {
	return &BookingHandler{
		store: s,
	}
}

// TODO: this requires ADMIN authorization
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// TODO: this requires USER authorization
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return c.JSON(booking)
}
