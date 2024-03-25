package api

import (
	"errors"
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(s *db.Store) *HotelHandler {
	return &HotelHandler{
		store: s,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {

	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
func (h *HotelHandler) HandlerGetHotels(c *fiber.Ctx) error {

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)

	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {

	id := c.Params("id")

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "document not found"})
		}

		return err
	}

	return c.JSON(hotel)
}
