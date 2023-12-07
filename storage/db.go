package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	return oid, err
}

const (
	DBURI = "mongodb://root:example@localhost:27017/"

	DBNAME    = "hotel-rezervation"
	USERCOLL  = "users"
	HOTELCOLL = "hotels"
	ROOMCOLL  = "rooms"
	BOOKCOLL = "book"

	TESTDBNAME = "test-hotel"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
	Booking BookingStore
}
