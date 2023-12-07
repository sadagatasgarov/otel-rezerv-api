package api

import (
	"fmt"
	"net/http"
	db "sadagatasgarov/hotel_rezerv_api/storage"
	"sadagatasgarov/hotel_rezerv_api/types"
	"time"

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
	fmt.Println(p.TillDate)
	if p.FromDate.Unix() > p.TillDate.Unix() {
		return fmt.Errorf("bitme tarixi(%+v) Baslama tarixi(%+v)-den kicik ola bilmez.", p.FromDate, p.TillDate)
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

func (h *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}
	//fmt.Println(params.FromDate)

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	//fmt.Println(roomID)
	user, ok := c.Context().Value("user").(*types.Users)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "Internel Server error",
		})
	}
	//fmt.Println(user)
	loc, _ := time.LoadLocation("Asia/Baku")
	where := bson.M{
		"roomID":roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate.In(loc),
		},
		"tillDate": bson.M{
			"$lte": params.TillDate.In(loc),
		},
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), where)
	if err != nil {
		return err
	}

	if len(bookings) > 0 {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
				Type: "error",
				Msg: fmt.Sprintf("room %+v already booked", roomID),
			})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate.In(loc),
		TillDate:   params.TillDate.In(loc),
		NumPersons: params.NumPersons,
	}

	//fmt.Println(booking)
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return fmt.Errorf("booking yazilmadi %v", err)
	}

	return c.JSON(inserted)
}
