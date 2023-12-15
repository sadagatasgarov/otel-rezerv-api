package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/middleware"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/storage/fixtures"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAdminGetBookings(t *testing.T) {
	var (
		tdb = setup(t)

		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "user", "admin", false)

		hotel   = fixtures.AddHotel(tdb.Store, "Test oteli", "Namelum yer", 5, nil)
		room    = fixtures.AddRoom(tdb.Store, "test _size", true, 50, hotel.ID, true)
		booking = fixtures.AddBooking(tdb.Store, user.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 2))

		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	defer tdb.teardown(t)
	// rezerv edilmis otaqlari siralayiriq
	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	req.Header.Add("Content-type", "application/json")
	getresp, err := app.Test(req, 3000)
	if err != nil {
		t.Fatal(err)
	}

	if getresp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", getresp.StatusCode)
	}

	bookings, err := tdb.Booking.GetBookings(context.TODO(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	var getbookings []*types.Booking
	if err := json.NewDecoder(getresp.Body).Decode(&getbookings); err != nil {
		t.Fatal(err)
	}

	if len(getbookings) != 1 {
		t.Fatalf("expected 1 booking gor %d ", len(getbookings))
	}

	if booking.ID != getbookings[0].ID {
		t.Fatal("Id ller uygun gelmir")
	}
	if !reflect.DeepEqual(bookings, getbookings) {
		fmt.Println(bookings)
		fmt.Println(getbookings)
		t.Fatal("expected bookinng to be equal")
	}

}
