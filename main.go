package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shoaibshazid/hotel-backend/api"
	"github.com/shoaibshazid/hotel-backend/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	clientOptions := options.Client().ApplyURI(db.DBURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	//handlers initialization
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)
		app          = fiber.New()
		apiV1        = app.Group("api/v1")
	)

	//user handlers
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandleUpdateUser)

	//hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	err = app.Listen(":6543")
	if err != nil {
		panic(err)
	}
}
