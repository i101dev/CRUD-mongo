package api

import (
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

type ResourceResp struct {
	Results int `json:"results"`
	Page    int `json:"page"`
	Data    any `json:"data"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
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
		return ErrorInvalidID()
	}

	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)

	if err != nil {
		return ErrorResourceNotFound("hotel data")
	}

	return c.JSON(rooms)
}
func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	var params HotelQueryParams

	if err := c.QueryParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := db.Map{
		"rating": params.Rating,
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)

	if err != nil {
		return ErrorResourceNotFound("all hotels data")
	}

	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
	}

	return c.JSON(resp)
}
func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {

	id := c.Params("id")

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)

	if err != nil {
		return ErrorResourceNotFound("hotel data")
	}

	return c.JSON(hotel)
}
