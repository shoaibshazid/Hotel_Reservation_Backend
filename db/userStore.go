package db

import (
	"context"
	"github.com/shoaibshazid/hotel-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollection = "users"

type Map map[string]any

type Dropper interface {
	Drop(ctx context.Context) error
}
type UserStore interface {
	Dropper
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error
	DeleteUser(ctx context.Context, userId string) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(UserCollection),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	curr, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := curr.All(ctx, &users); err != nil {
		return []*types.User{}, err
	}
	return users, err
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, userId string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	// TODO: Maybe its a good ides to handle if we did not delete the user
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}
	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
