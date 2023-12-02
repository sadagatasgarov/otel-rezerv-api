package main

import (
	"context"
	"flag"
	"log"

	"sadagatasgarov/hotel_rezerv_api/api"
	db "sadagatasgarov/hotel_rezerv_api/storage"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "The listen addres of the API server")
	flag.Parse()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore   = db.NewMongoUserStore(client)
		userHandler = api.NewUserHandler(userStore)

		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		store      = db.Store{
			Hotel: hotelStore,
			User:  userStore,
			Room:  roomStore,
		}
		hotelHandler = api.NewHotelHandler(&store)

		app   = fiber.New(config)
		apiv1 = app.Group("api/v1")
	)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)

	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listenAddr)
}

// app.Get("/foo", hanlerFunc)
// func hanlerFunc(c *fiber.Ctx) error {
// 	return c.JSON(map[string]string{"msg": "working"})
// }
