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
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
	app := fiber.New()
	apiV1 := app.Group("api/v1")
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandleUpdateUser)
	err = app.Listen(":6543")
	if err != nil {
		panic(err)
	}
}
