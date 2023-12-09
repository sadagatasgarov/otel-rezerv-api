package main

import (
	"context"
	"fmt"
	"log"
	db "sadagatasgarov/hotel_rezerv_api/storage"
	"sadagatasgarov/hotel_rezerv_api/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)



func seedUser(isadmin bool, fname, lname, email, password string) {
	params := types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
		IsAdmin:   isadmin,
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		fmt.Println(err)
	}

	use, err := userStore.InsertUser(ctx, user)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(use)

}

func seedHotel(name string, location string, rating int) {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []*primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 88.9,
			//HotelID: hotel.ID,
		},
		{
			Size: "normal",
			Price: 100,
			//HotelID: hotel.ID,
		},
		{
			Size:  "kingsize",
			Price: 200,
			//HotelID: hotel.ID,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	seedHotel("Belus", "franc", 3)
	seedHotel("Cozy", "franc", 2)
	seedHotel("Hinter", "USA", 5)
	seedUser(true, "admin", "admin", "admin@admin.com", "adminadmin")
	seedUser(false, "user", "user", "user@user.com", "useruser")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
}
