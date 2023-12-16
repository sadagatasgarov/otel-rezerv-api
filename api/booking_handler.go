package api

import (
	"net/http"

	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(
			genericResp{
				Type: "info",
				Msg:  "not authorized",
			})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{"cancelled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "Cancelled",
	})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrNotResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(
			genericResp{
				Type: "info",
				Msg:  "Bilinmeyen booking id",
			})

	}
	return c.JSON(booking)
}
