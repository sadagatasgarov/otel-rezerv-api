package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sadagatasgarov/otel-rezerv-api/api"
	db "gitlab.com/sadagatasgarov/otel-rezerv-api/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	//listenAddr := flag.String("listenAddr", ":5000", "The listen addres of the API server")

	//flag.Parse()
	var client *mongo.Client
	var err error
	//client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	
	if db.DBURI == ""{
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURILOKAL))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI).SetServerAPIOptions(serverAPI))
		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		userHandler  = api.NewUserHandler(userStore)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookStore(client)
		store        = db.Store{
			Hotel:   hotelStore,
			User:    userStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(&store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(&store)
		bookingHandler = api.NewBookingHandler(&store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	// Versioned API routes
	auth.Post("/auth", authHandler.HandleAuth)
	auth.Post("/user", userHandler.HandleCreateUser)

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

	// rooms handlers
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	// bookings handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Post("/booking/:id", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTENING_PORT")
	if listenAddr == "" {
		listenAddr := ":6000"	
		app.Listen(listenAddr)
	}else{
		app.Listen(listenAddr)
	}
	

}

// app.Get("/foo", hanlerFunc)
// func hanlerFunc(c *fiber.Ctx) error {
// 	return c.JSON(map[string]string{"msg": "working"})
// }

//39  18:49
