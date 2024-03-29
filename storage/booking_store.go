package db

import (
	"context"
	"fmt"

	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	BookingStore
}

func NewMongoBookStore(client *mongo.Client) *MongoBookingStore {
	if DBNAME == "" {
		return &MongoBookingStore{
			client: client,
			coll:   client.Database(DBNAMELOKAL).Collection(BOOKCOLL),
		}
	} else {
		return &MongoBookingStore{
			client: client,
			coll:   client.Database(DBNAME).Collection(BOOKCOLL),
		}
	}

}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping user collection bu isledi")
	return s.coll.Drop(ctx)
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update = bson.M{"$set": update}
	resp, err := s.coll.UpdateByID(ctx, oid, update)
	_ = resp

	return err
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {

	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking

	if err := resp.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking *types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}
