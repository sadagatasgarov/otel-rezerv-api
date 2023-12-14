package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"gitlab.com/sadagatasgarov/otel-rezerv-api/types"

	"testing"

	"github.com/gofiber/fiber/v2"
)

//opts := options.Client().ApplyURI("mongodb+srv://<username>:<password>@cluster0.nlvrqpz.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

//var testuri = "mongodb://root:example@localhost:27017/"
//var testuriatlas = "mongodb+srv://user:example@cluster0.nlvrqpz.mongodb.net/?retryWrites=true&w=majority"
//mongodb+srv://<username>:<password>@cluster0.nlvrqpz.mongodb.net/?retryWrites=true&w=majority
//var testdbname = "test-hotel-rezervation"
// var testuserColl = "test-users"

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()

	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "user2",
		LastName:  "user2",
		Email:     "user2@user.com",
		Password:  "useruser",
		IsAdmin:   false,
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req, 3000)
	if err != nil {
		t.Fatal(err)
	}

	var user types.Users
	//bb, _ := io.ReadAll(resp.Body)
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Fatalf(" %s User id 0 ola bilmezx", user.ID)
	}

	if len(user.EncryptedPassword) > 0 {
		t.Fatalf(" %s EncryptedPassword not to be in the json response", user.EncryptedPassword)
	}

	if user.FirstName != params.FirstName {
		t.Fatalf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Fatalf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Fatalf("expected Email %s but got %s", params.Email, user.Email)
	}

	//fmt.Println(user)
}
