package users

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User is the schema for the user model
type User struct {
	ID           int64       `json:"id"`
	Email        string      `json:"email"`
	FName        string      `json:"first_name"`
	LName        string      `json:"last_name"`
	Role         string      `json:"role"`
	DatesPresent []time.Time `json:"dates_present"`
	Password     string      `json:"password"`
}

// UserPassChange is the struct for changing the password of a user
type UserPassChange struct {
	Password string `json:"password"`
}

// UserLoginStruct is the struct for logging a user in
type UserLoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWTToken is the struct for the JWT token
type JWTToken struct {
	Email string `json:"email"`
	FName string `json:"first_name"`
	LName string `json:"last_name"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (u User) String() string {
	return fmt.Sprintf("%d - %s %s - %s", u.ID, u.FName, u.LName, u.Email)
}

// GetDatesPresent returns the number of days present
func (u User) GetDatesPresent() []time.Time {
	return u.DatesPresent
}
