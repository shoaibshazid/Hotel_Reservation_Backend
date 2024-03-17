package main

import (
	"context"
	"fmt"
	"github.com/shoaibshazid/hotel-backend/db"
	"github.com/shoaibshazid/hotel-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(email string, firstName string, lastName string) error {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  "Password@123",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func seedHotel(name string, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: 99.9,
		},
		{
			Size:      "normal",
			BasePrice: 999.9,
		},
		{
			Size:      "large",
			BasePrice: 9999.9,
		},
		{
			Size:      "King-Size",
			BasePrice: 99999.9,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeding the database")
	fmt.Println(insertedHotel)
	for _, room := range rooms {
		room.HotelId = insertedHotel.Id
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
	return nil
}

func main() {
	var err error
	err = seedHotel("Meridian", "Hyderabad", 5)
	if err != nil {
		return
	}
	err = seedHotel("Sheraton", "Hyderabad", 4)
	if err != nil {
		return
	}
	err = seedUser("shoaib1123@gmail.com", "shoaib", "sajid")
}

func init() {
	var err error
	clientOptions := options.Client().ApplyURI(db.DBURI)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	fmt.Println(client)
}
