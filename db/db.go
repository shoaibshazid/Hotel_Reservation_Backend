package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const Dbname = "hotel-booking"

func ToObjectId(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid, err
}
