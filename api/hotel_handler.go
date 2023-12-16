package api

import (
	"fmt"

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
	Results int `json:"results"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	Rating int `json:"rating"`
	db.Pagination
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	//var pagination db.Pagination
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}


	filter := db.Map{
		"rating": params.Rating,
	}
	//fmt.Println(params.Rating, params.Page, params.Limit)

	if params.Page==0{
		params.Limit=0
	}

	if params.Rating==0{
		filter=nil
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrNotResourceNotFound("hotels")
	}

	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
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
