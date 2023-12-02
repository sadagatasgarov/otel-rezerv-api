package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	db "sadagatasgarov/hotel_rezerv_api/storage"
	"sadagatasgarov/hotel_rezerv_api/types"

	"log"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testuri = "mongodb://root:example@localhost:27017/"

var testdbname = "test-hotel-rezervation"

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
	app := fiber.New()

	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "test Ad",
		LastName:  "test Soyad",
		Email:     "test@test.dsd",
		Password:  "12345678",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.Users
	//bb, _ := io.ReadAll(resp.Body)
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf(" %s User id 0 ola bilmezx", user.ID)
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf(" %s EncryptedPassword not to be in the json response", user.EncryptedPassword)
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected Email %s but got %s", params.Email, user.Email)
	}

	//fmt.Println(user)
}
