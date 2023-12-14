package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.com/sadagatasgarov/otel-rezervasiya-api/api"
	db "gitlab.com/sadagatasgarov/otel-rezervasiya-api/storage"
	"gitlab.com/sadagatasgarov/otel-rezervasiya-api/storage/fixtures"

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
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Booking: db.NewMongoBookStore(client),
	}

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println(admin)
	fmt.Println("user -->", admin.ID, api.CreateTokenFromUser(admin))

	user := fixtures.AddUser(store, "sada", "asga", false)
	fmt.Println(user)
	fmt.Println("user -->", user.ID, api.CreateTokenFromUser(user))

	hotel := fixtures.AddHotel(store, "Tebriz", "Nakhchivan", 5, nil)
	fmt.Println("hotel ----> ", hotel)

	room := fixtures.AddRoom(store, "small", true, 110.55, hotel.ID, true)
	fmt.Println("room ----> ", room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ----> ", booking)

}
