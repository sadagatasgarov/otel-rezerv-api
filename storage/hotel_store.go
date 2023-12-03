package db

import (
	"context"
	"sadagatasgarov/hotel_rezerv_api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filter bson.M, update bson.M) error
	GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error)
	GetHotel(ctx context.Context, oid primitive.ObjectID) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(HOTELCOLL),
	}
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) GetHotel(ctx context.Context, oid primitive.ObjectID) (*types.Hotel, error) {
	var hotel *types.Hotel
	if err:=s.coll.FindOne(ctx, bson.M{"_id":oid}).Decode(&hotel); err!=nil{
		return nil, err
	}
	return hotel, nil
}
