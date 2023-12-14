package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	//fmt.Println(p.TillDate)
	if p.FromDate.Unix() > p.TillDate.Unix() {
		return fmt.Errorf("bitme tarixi(%+v) Baslama tarixi(%+v)-den kicik ola bilmez", p.FromDate, p.TillDate)
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}
	fmt.Println(params.FromDate)

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.Users)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internel Server error",
		})
	}

	//fmt.Println(roomID, params)
	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)

	if err != nil {
		return fmt.Errorf("aaaaaaaaaa")
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %+v already booked", roomID),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	//fmt.Println(booking)
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return fmt.Errorf("booking yazilmadi %v", err)
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, nil
	}

	ok := len(bookings) == 0
	if !ok {
		return false, nil
	}
	return ok, nil
}
