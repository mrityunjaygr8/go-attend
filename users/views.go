package users

import (
	"errors"

	"github.com/go-pg/pg"
)

// GetUsers fetches all users
func GetUsers(db *pg.DB) ([]User, error) {
	var users []User
	err := db.Model(&users).Order("id ASC").Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser creates a new user in the database
func (u *User) CreateUser(db *pg.DB) error {
	_, err := db.Model(u).Insert()
	if err != nil {
		return err
	}
	return nil
}

// GetUser gets a user from the database using the given id
func (u *User) GetUser(db *pg.DB) (*User, error) {
	count, err := db.Model(u).WherePK().Count()

	if count < 1 {
		return nil, errors.New("User does not exist")
	}
	err = db.Model(u).WherePK().Select()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteUser deletes a user from the database
func (u *User) DeleteUser(db *pg.DB) error {
	count, err := db.Model(u).WherePK().Count()

	if count < 1 {
		return errors.New("User does not exist")
	}
	_, err = db.Model(u).WherePK().Delete()

	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates the user
func (u *User) UpdateUser(db *pg.DB) error {
	count, err := db.Model(u).WherePK().Count()

	if count < 1 {
		return errors.New("User does not exist")
	}
	_, err = db.Model(u).WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}
