package api

import (
	"context"
	db "sadagatasgarov/hotel_rezerv_api/storage"

	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testuri = "mongodb://root:example@localhost:27017/"

// var testdbname = "test-hotel-rezervation"
// var testuserColl = "test-users"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {

	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}

}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testuri))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
}
