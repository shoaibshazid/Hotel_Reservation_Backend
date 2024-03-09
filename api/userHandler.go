package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/shoaibshazid/hotel-backend/db"
	"github.com/shoaibshazid/hotel-backend/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"data": "userId is deleted"})
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return nil
	}

	insertedUser, e := h.userStore.CreateUser(c.Context(), user)

	if e != nil {
		return e
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"message": "User not found"})
		}
		return err
	}
	return c.JSON(users)
}
