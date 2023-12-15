package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/storage/fixtures"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "admin", "admin", true)
	user := fixtures.AddUser(tdb.Store, "user", "admin", false)

	hotel := fixtures.AddHotel(tdb.Store, "Test oteli", "Namelum yer", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "test _size", true, 50, hotel.ID, true)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 2))

	params := AuthParams{
		Email:    "admin@admin.com",
		Password: "admin_admin",
	}
	b, _ := json.Marshal(params)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)

	app.Post("/bookpost", authHandler.HandleAuth)

	req := httptest.NewRequest(http.MethodPost, "/bookpost", bytes.NewReader(b))
	resp1, err := app.Test(req, 2000)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]string
	json.NewDecoder(resp1.Body).Decode(&data)
	token := data["token"]

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-token", token)



	bookHandler := NewBookingHandler(tdb.Store)
	app.Get("/booklist", bookHandler.HandleGetBookings)
	req = httptest.NewRequest(http.MethodGet, "/booklist", nil)
	resp2, err := app.Test(req, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp2.StatusCode)
	}

	// bookings, err := tdb.Booking.GetBookings(context.TODO(), bson.M{})
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var respbooking *types.Booking
	if err := json.NewDecoder(resp1.Body).Decode(&respbooking); err != nil {
		t.Fatal(err)
	}

	var respbookings []*types.Booking
	if err := json.NewDecoder(resp2.Body).Decode(&respbookings); err != nil {
		t.Fatal(err)
	}

	if len(respbookings) != 1 {
		t.Fatalf("expected 1 booking gor %d ", len(respbookings))
	}

	if !reflect.DeepEqual(respbooking, respbookings[0]) {
		fmt.Println(booking)
		fmt.Println(respbookings[0])
		t.Fatal("expected bookinng to be equal")
	}
}
