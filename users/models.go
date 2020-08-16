package users

import (
	"fmt"
	"time"
)

// User is the schema for the user model
type User struct {
	ID           int64       `json:"id"`
	Email        string      `json:"email"`
	FName        string      `json:"first_name"`
	LName        string      `json:"last_name"`
	Role         string      `json:"role"`
	DatesPresent []time.Time `json:"dates_present"`
}

func (u User) String() string {
	return fmt.Sprintf("%d - %s %s - %s", u.ID, u.FName, u.LName, u.Email)
}

// GetDatesPresent returns the number of days present
func (u User) GetDatesPresent() []time.Time {
	return u.DatesPresent
}
