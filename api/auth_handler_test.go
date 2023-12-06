package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	db "sadagatasgarov/hotel_rezerv_api/storage"
	"sadagatasgarov/hotel_rezerv_api/types"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func makeTestUser(t *testing.T, userStore db.UserStore) *types.Users {
	params := types.CreateUserParams{
		FirstName: "fname",
		LastName:  "lname",
		Email:     "email@tst.txt",
		Password:  "12345678",
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Fatal(err)
	}

	use, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return use
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "email@tst.txt",
		Password: "12345678",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatalf("expected the user to be inserted user")
	}
}

func TestAuthenticateWrongPassFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "email@tst.txt",
		Password: "123456789dsda",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400 but %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
}
