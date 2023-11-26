package db

import (
	"context"
	"hotel_api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBNAME = "hotel-rezervation"
const USERCOLL = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.Users, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {

	return &MongoUserStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(USERCOLL),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.Users, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.Users
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
