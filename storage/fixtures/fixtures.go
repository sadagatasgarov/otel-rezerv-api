package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	db "gitlab.com/sadagatasgarov/otel-rezervasiya-api/storage"
	"gitlab.com/sadagatasgarov/otel-rezervasiya-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, uid primitive.ObjectID, rid primitive.ObjectID, numPersons int, from time.Time, till time.Time) *types.Booking {

	booking := types.Booking{
		UserID:     uid,
		RoomID:     rid,
		NumPersons: numPersons,
		FromDate:   from,
		TillDate:   till,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID, available bool) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []*primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []*primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddUser(store *db.Store, fname, lname string, isadmin bool) *types.Users {
	params := types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
		IsAdmin:   isadmin,
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}
