package api

import (
	"hotel-reservation/db"
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

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "nost authorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Msg: "booking cancelled"})

}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "nost authorized",
		})
	}

	return c.JSON(booking)
}
