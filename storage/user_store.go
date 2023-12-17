package db

import (
	"context"
	"fmt"

	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Map map[string]any

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper
	GetUserById(context.Context, string) (*types.Users, error)
	GetUsers(context.Context) ([]*types.Users, error)
	InsertUser(context.Context, *types.Users) (*types.Users, error)
	DeleteUser(context.Context, string) (*types.Users, error)
	UpdateUser(context.Context, string, types.UpdateUserParams) error
	GetUserByEmail(context.Context, string) (*types.Users, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {

	if DBNAME == "" {
		return &MongoUserStore{
			client: c,
			coll:   c.Database(DBNAMELOKAL).Collection(USERCOLL),
		}
	} else {
		return &MongoUserStore{
			client: c,
			coll:   c.Database(DBNAME).Collection(USERCOLL),
		}
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping user collection bu isledi")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := Map{"_id": oid}
	update := Map{"$set": params.ToBSON()}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) (*types.Users, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.Users
	if err := s.coll.FindOne(ctx, Map{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	_, err = s.coll.DeleteOne(ctx, Map{"_id": user.ID})

	return &user, err
}

func (s *MongoUserStore) InsertUser(ctx context.Context, u *types.Users) (*types.Users, error) {

	// var user *types.Users
	// if err := s.coll.FindOne(ctx, Map{"email": u.Email}).Decode(&user); err != nil {
	// 	res, err := s.coll.InsertOne(ctx, u)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	u.ID = res.InsertedID.(primitive.ObjectID)
	// 	return u, nil
	// }
	// if reflect.DeepEqual(user.Email, u.Email) {
	// 	return nil, fmt.Errorf("email bazada movcuddur")
	// }

	// return u, nil

	res, err := s.coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
	return u, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.Users, error) {
	var users []*types.Users
	cur, err := s.coll.Find(ctx, Map{})
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
	if err := s.coll.FindOne(ctx, Map{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.Users, error) {
	var user types.Users
	if err := s.coll.FindOne(ctx, Map{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
