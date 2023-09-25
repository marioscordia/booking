package main

import (
	"booking/api"
	"booking/api/middleware"
	"booking/config"
	"booking/db"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new fiber instance with custom config
var (
	conf = fiber.Config{
	// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			
			if e, ok := err.(api.Error); ok {
				return ctx.Status(e.Code).JSON(map[string]string{"msg": e.Msg,})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(map[string]string{
				"msg": http.StatusText(http.StatusInternalServerError),
			})
		},
	}
)

func main() {

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DB.URI))
	if err != nil {
		log.Fatal(err)
	}

	DB := db.NewDB(client, config)

	userHandler := api.NewUserHandler(DB)
	hotelHandler := api.NewHotelHandler(DB)
	authHandler := api.NewAuthHandler(DB)
	roomHandler := api.NewRoomHandler(DB)

	app := fiber.New(conf)
	user := app.Group("/user")
	admin := app.Group("/admin", middleware.JWTauthentication, middleware.IsAdmin)
	apiBook := app.Group("/room", middleware.JWTauthentication)

	// Admin handlers
	admin.Get("/user", userHandler.HandleGetUsers)
	admin.Get("/user/:id", userHandler.HandleGetUser)
	admin.Delete("/user/:id", userHandler.HandleDeleteUser)

	//User handlers
	user.Post("/register", userHandler.HandlePostUser)
	user.Put("/update", userHandler.HandlePutUser, middleware.JWTauthentication)
	user.Post("/login", authHandler.Login)

	//Hotel handlers
	app.Get("/hotel", hotelHandler.HandleGetHotels)
	app.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByID)
	
	//Booking handlers
	apiBook.Post("/:roomid", roomHandler.HandleBookRoom)
	apiBook.Get("/booked", roomHandler.HandleGetBookings)
	apiBook.Delete("/:id", roomHandler.HandleCancelBooking)
	
	fmt.Printf("Working on localhost%s\n", config.Port)
	log.Fatal(app.Listen(config.Port))
	
}
