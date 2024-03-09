package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptcost      = 12
	minFirstNameLen = 7
	minLastNameLen  = 7
	minPasswordLen  = 8
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"first_name" json:"firstName"`
	LastName          string             `bson:"last_name" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encrypted_password" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptcost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPassword),
	}, nil
}

func (params CreateUserParams) Validate() error {
	if len(params.FirstName) < minFirstNameLen {
		return fmt.Errorf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		return fmt.Errorf("lastname length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		return fmt.Errorf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
