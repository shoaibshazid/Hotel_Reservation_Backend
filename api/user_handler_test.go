package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shoaibshazid/hotel-backend/db"
	"github.com/shoaibshazid/hotel-backend/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"net/http/httptest"
	"testing"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	clientOptions := options.Client().ApplyURI(db.TESTDBURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.DBNAMETESTING),
	}
}
func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	UserHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", UserHandler.HandlePostUser)
	params := types.CreateUserParams{
		Email:     "hello@gmail.com",
		FirstName: "Cilian Murphy",
		LastName:  "Thomas Shelby",
		Password:  "ihdh6689dbhhhshc",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	response, _ := app.Test(req)
	fmt.Println(response.Status)
	var user types.User
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected lastname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}
