package main

import (
	"booking/config"
	"booking/db"
	"booking/types"
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	wg sync.WaitGroup
)

func main() {
	fmt.Println("Creating a database...")

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	if err := client.Database(config.DB.Name).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, config)
	roomStore := db.NewMongoRoomStore(client, config)
	userStore := db.NewMongoUserStore(client, config)
	for _, h := range db.Hotels {
		wg.Add(1)
		go func (hotel types.Hotel, rooms []types.Room) {  
			insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
			if err != nil{
				log.Fatal(err)
			}
			
			for _, r := range rooms{
				r.HotelID = insertedHotel.ID
			
				room, err := roomStore.InsertRoom(ctx, &r)
				if err != nil{
					log.Fatal(err)
				}
		
				filter, update := bson.M{"_id": room.HotelID}, bson.M{"$push":bson.M{"rooms": room.ID}}
		
				err = hotelStore.Update(ctx, filter, update)
				if err != nil{
					log.Fatal(err)
				}
			}
			wg.Done()
		}(h, db.Rooms)
	}

	wg.Wait()

	us := client.Database(config.DB.Name).Collection(config.DB.User)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // 1 for ascending, -1 for descending
		Options: options.Index().SetUnique(true),
	}
	_, err = us.Indexes().CreateOne(ctx, indexModel)
	if err != nil{
		log.Fatal(err)
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte("Cheburek"), 12)
	if err != nil {
		log.Fatal(err)
	}

	admin := types.User{
		FirstName: "James",
		LastName: "Bond",
		Email: "007@gmail.com",
		IsAdmin: true,
		Password: string(encpw),
	}
	err = userStore.InsertUser(ctx, &admin)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Database created...")
}

