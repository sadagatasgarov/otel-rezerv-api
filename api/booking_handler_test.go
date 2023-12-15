package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/storage/fixtures"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "admin", "admin", true)
	user := fixtures.AddUser(tdb.Store, "user", "admin", false)


	hotel := fixtures.AddHotel(tdb.Store, "Test oteli", "Namelum yer", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "test _size", true, 50, hotel.ID, true)
	fixtures.AddBooking(tdb.Store, user.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 2))
	
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/booklist", authHandler.HandleAuth)
	params := AuthParams{
		Email:    "admin@admin.com",
		Password: "admin_admin",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/booklist", bytes.NewReader(b))
	resp, err := app.Test(req, 2000)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	token := data["token"]

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-token", token)

	listbooking, err := tdb.Booking.GetBookings(context.TODO(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range listbooking {
		fmt.Println(v)

	}

	//fmt.Println(CreateTokenFromUser(insertedUser))
}
