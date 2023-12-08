package api

import (
	db "sadagatasgarov/hotel_rezerv_api/storage"

	"github.com/gofiber/fiber/v2"
)



type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler{
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetRooms(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}
