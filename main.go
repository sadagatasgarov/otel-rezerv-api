package main

import (
	"context"
	"flag"
	"hotel_api/api"
	db "hotel_api/storage"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri = "mongodb://root:example@localhost:27017/"
var dbname = "hotel-rezervation"
var userColl = "users"

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
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}

func hanlerFunc(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working"})
}

// coll := client.Database(dbname).Collection(userColl)

// user := types.Users{
// 	FirstName: "Sada",
// 	LastName:  "At the water cooler",
// }

// res, err := coll.InsertOne(context.Background(), user)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println(res)

// var james types.Users
// if err:=coll.FindOne(context.Background(), bson.M{}).Decode(&james); err!=nil{
// 	log.Fatal(err)
// }
// fmt.Println(james)
