package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 3
	minLastNameLen  = 3
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {

	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
}
func (p UpdateUserParams) ToBSON() bson.M {

	m := bson.M{}

	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	return m
}
func (params CreateUserParams) Validate() map[string]string {

	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstname minimum length: %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName minimum length: %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password minimum length: %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = "invalid email"
	}

	return errors
}
