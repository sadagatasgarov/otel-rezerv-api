package main

import (
	"context"
	"flag"
	"sadagatasgarov/hotel_rezerv_api/api"
	db "sadagatasgarov/hotel_rezerv_api/storage"

	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri = "mongodb://root:example@localhost:27017/"

//var dbname = "hotel-rezervation"
//var userColl = "users"

// Create a new fiber instance with custom config
var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "The listen addres of the API server")
	flag.Parse()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")

	app.Get("/foo", hanlerFunc)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)

	app.Listen(*listenAddr)
}

func hanlerFunc(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working"})
}
