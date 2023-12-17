package db

import "os"

var DBURI string = os.Getenv("MONGODB_URL")
var DBNAME string = os.Getenv("DB_NAME")
//var DBURIATLAS string = os.Getenv("MONGODB_URL_TEST")

const (
	DBURILOKAL = "mongodb://root:example@localhost:27017/"
	//DBURIATLAS = "mongodb+srv://user:example@cluster0.nlvrqpz.mongodb.net/?retryWrites=true&w=majority"
	DBNAMELOKAL     = "hotel-rezervation"
	USERCOLL  = "users"
	HOTELCOLL = "hotels"
	ROOMCOLL  = "rooms"
	BOOKCOLL  = "book"
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
