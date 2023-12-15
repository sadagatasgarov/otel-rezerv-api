package api

import (
	"context"
	"log"
	"testing"

	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {

	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}

}

func setup(t *testing.T) *testdb {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURIATLAS).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			User:    db.NewMongoUserStore(client),
			Booking: db.NewMongoBookStore(client),
		},
	}
}
