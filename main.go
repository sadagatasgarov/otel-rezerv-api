package main

import (
	"context"
	"flag"
	"log"

	"sadagatasgarov/hotel_rezerv_api/api"
	"sadagatasgarov/hotel_rezerv_api/middleware"
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
		userStore   = db.NewMongoUserStore(client, db.DBNAME)
		userHandler = api.NewUserHandler(userStore)
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		store       = db.Store{
			Hotel: hotelStore,
			User:  userStore,
			Room:  roomStore,
		}
		hotelHandler = api.NewHotelHandler(&store)
		
		authHandler = api.NewAuthHandler(userStore)

		roomHandler = api.NewRoomHandler(&store)
		app   = fiber.New(config)
		api   = app.Group("api")
		apiv1 = app.Group("api/v1", middleware.JWTAuthentication(userStore))
	)


	// Versioned API routes
	api.Post("/auth", authHandler.HandleAuth)

	// user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Post("/room/:id/book", roomHandler.HandleRookRoom)
	app.Listen(*listenAddr)
}

// app.Get("/foo", hanlerFunc)
// func hanlerFunc(c *fiber.Ctx) error {
// 	return c.JSON(map[string]string{"msg": "working"})
// }
