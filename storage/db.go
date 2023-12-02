package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	return oid, err
}

const (
	DBNAME   = "hotel-rezervation"
	USERCOLL = "users"
	DBURI    = "mongodb://root:example@localhost:27017/"
)

type Store struct {
	User UserStore
	Hotel HotelStore
	Room RoomStore
}
