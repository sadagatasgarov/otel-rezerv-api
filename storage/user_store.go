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
	GetUsers(context.Context) ([]*types.Users, error)
	InsertUser(context.Context, *types.Users) (*types.Users, error)
	DeleteUser(context.Context, string) (*types.Users, error)
	UpdateUser(ctx context.Context, filter, update bson.M) (error)
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

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter, values bson.M) error {
	update := bson.D{
		{
			Key: "$set", Value: values,
		},
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, u *types.Users) (*types.Users, error) {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var user types.Users
// 	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
// 		return nil, err
// 	}

// 	_, err = s.coll.UpdateOne(ctx, user, u)

// 	return &user, nil
// }

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) (*types.Users, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.Users
	if err := s.coll.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&user); err != nil {
		return nil, err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": user.ID})

	return &user, err
}

func (s *MongoUserStore) InsertUser(ctx context.Context, u *types.Users) (*types.Users, error) {
	// usr, err := types.NewUserFromParams()

	res, err := s.coll.InsertOne(ctx, u)

	if err != nil {
		return nil, err
	}

	u.ID = res.InsertedID.(primitive.ObjectID)

	return u, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.Users, error) {
	var users []*types.Users

	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
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
