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
	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "admin", "admin", true)

	//hotel := fixtures.AddHotel(tdb.Store, "Test oteli", "Namelum yer", 5, nil)
	//room := fixtures.AddRoom(tdb.Store, "test _size", true, 50, hotel.ID, true)
	//fixtures.AddBooking(tdb.Store, user1.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 2))
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "sada",
		LastName:  "asga",
		Email:     "sada@asga.com",
		Password:  "sada_asga",
		IsAdmin:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
	user, err = tdb.User.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	
	var roomIDS []*primitive.ObjectID
	if roomIDS == nil {
		roomIDS = []*primitive.ObjectID{}
	}
	hotel, err := tdb.Hotel.Insert(context.TODO(), &types.Hotel{
		Name:     "Test Oteli",
		Location: "namelum_yer",
		Rooms:    roomIDS,
		Rating:   5,
	})
	if err != nil {
		t.Fatal(err)
	}

	room, err := tdb.Room.InsertRoom(context.TODO(), &types.Room{
		Size:      "Luks",
		Seaside:   true,
		Price:     55,
		HotelID:   hotel.ID,
		Available: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	booking, err := tdb.Booking.InsertBooking(context.TODO(), &types.Booking{
		UserID:     user.ID,
		RoomID:     room.ID,
		NumPersons: 5,
		FromDate:   time.Now(),
		TillDate:   time.Now().AddDate(0, 0, 2),
	})
	if err != nil {
		t.Fatal(err)
	}

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

	fmt.Println(booking.ID)
	for _, v := range listbooking {
		fmt.Println(v)
	}

}
