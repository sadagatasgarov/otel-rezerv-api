package api

import (
	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	//fmt.Println(&oid)
	filter := db.Map{"hotelID": &oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNotResourceNotFound("rooms")
	}
	return c.JSON(rooms)
}

type ResourceResp struct {
	Data    any `json:"data"`
	Results int   `json:"results"`
	Page    int   `json:"page"`
}

// func NewResourceResp(data any, page int) ResourceResp {
// 	return ResourceResp{
// 		Results: len(data),
// 		Data: data,
// 		Page: page,
// 	}
// }

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var pagination db.Pagination
	if err := c.QueryParser(&pagination); err != nil {
		return ErrBadRequest()
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &pagination)
	if err != nil {
		return ErrNotResourceNotFound("hotels")
	}

	resp:=ResourceResp{
		Data: hotels,
		Results: len(hotels),
		Page: int(pagination.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotel(c.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}
