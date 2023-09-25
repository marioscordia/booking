package db

import (
	"booking/config"
	"booking/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	UserStore
	HotelStore
	RoomStore
	BookStore
}

func NewDB(client *mongo.Client,  config *config.Config) *DB{
	return &DB{
		NewMongoUserStore(client, config),
		NewMongoHotelStore(client, config),
		NewMongoRoomStore(client, config),
		NewMongoBookStore(client, config),
	}
}

// const (
// 	DBuri = "mongodb://localhost:27017"
// 	DBname = "hotel-reservation"
// 	Usercoll = "users"
// 	Hotelcoll = "hotels"
// 	Roomcoll = "rooms"
// 	Bookcoll = "bookings"
// )

var (
	Hotels = []types.Hotel{
		{
			Name: "Meder",
			Location: "Kemer",
			Rooms: []primitive.ObjectID{},
		},
		{
			Name: "Hayat",
			Location: "Antalya",
			Rooms: []primitive.ObjectID{},
		},
		{
			Name: "Akpa",
			Location: "Bodrum",
			Rooms: []primitive.ObjectID{},
		},
	}
	Rooms = []types.Room{
		{
			Type: types.Single,
			Price: 99.9,
			Bed: types.Small,
			SeaSide: false,
		},
		{
			Type: types.Double,
			Price: 149.9,
			Bed: types.Medium,
			SeaSide: false,
		},
		{
			Type: types.Deluxe,
			Price: 299.9,
			Bed: types.Large,
			SeaSide: true,
		},
	}
)

