package users

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
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
	count, err := db.Model(u).Where("email = ?", u.Email).Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("User already exists")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	u.Password = string(pass)
	_, err = db.Model(u).Insert()
	if err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the password for a user
func (u *User) UpdatePassword(db *pg.DB) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	u.Password = string(pass)
	_, err = db.Model(u).WherePK().Update()
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

// GetForLogin gets the user for login
func GetForLogin(db *pg.DB, email string) (*User, error) {
	user := &User{}
	count, err := db.Model(user).Where("email = ?", email).Count()
	if err != nil {
		return nil, err
	}
	if count < 1 {
		return nil, errors.New("User does not exist")
	}

	err = db.Model(user).Where("email = ?", email).Select()

	return user, nil
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
