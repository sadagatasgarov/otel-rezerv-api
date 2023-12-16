package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"gitlab.com/sadagatasgarov/otel-rezerv-api/api"
	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/storage/fixtures"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookStore(client),
	}

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin -->", admin.Email, admin.ID, api.CreateTokenFromUser(admin))
	fmt.Println()
	user := fixtures.AddUser(store, "user", "user", false)
	fmt.Println("user  -->", user.Email, user.ID, api.CreateTokenFromUser(user))
	fmt.Println()

	hotel := fixtures.AddHotel(store, "Tebriz", "Nakhchivan", 5, nil)
	fmt.Println("hotel ----> ", hotel)

	room := fixtures.AddRoom(store, "small", true, 110.55, hotel.ID, true)
	fmt.Println("room ----> ", room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ----> ", booking)

	for i := 0; i < 10; i++ {

		fixtures.AddHotel(store, "Tebriz", "Nakhchivan", 5, nil)

	}
	for i := 0; i < 5; i++ {

		fixtures.AddHotel(store, "Tebriz", "Nakhchivan", 3, nil)

	}

	for i := 0; i < 20; i++ {

		fixtures.AddHotel(store, "Tebriz", strconv.Itoa(i), 2, nil)

	}

}
