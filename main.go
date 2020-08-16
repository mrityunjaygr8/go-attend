package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mrityunjaygr8/go-attend/app"
	"github.com/mrityunjaygr8/go-attend/users"
)

func main() {
	godotenv.Load()
	user := &users.User{
		Email:        "user@example.com",
		FName:        "user",
		LName:        "example",
		DatesPresent: []time.Time{time.Now(), time.Now().AddDate(0, 0, -3)},
		Role:         "admin",
	}

	newUser := &users.User{
		Email:        "user@example.com",
		FName:        "user",
		LName:        "example",
		DatesPresent: []time.Time{time.Now(), time.Now().AddDate(0, 0, -3)},
		Role:         "base",
	}

	app := app.App{Port: ":8000"}
	config := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("ATTEND_DATABASE_TYPE"),
		os.Getenv("ATTEND_DATABASE_USER"),
		os.Getenv("ATTEND_DATABASE_PASS"),
		os.Getenv("ATTEND_DATABASE_HOST"),
		os.Getenv("ATTEND_DATABASE_PORT"),
		os.Getenv("ATTEND_DATABASE_NAME"),
	)
	err := app.Setup(config)
	defer app.Close()
	if err != nil {
		panic(err)
	}

	db := app.Db

	_, err = db.Model(user, newUser).Insert()

	user2 := &users.User{
		ID: user.ID,
	}

	err = db.Model(user2).WherePK().Select()
	if err != nil {
		panic(err)
	}

	user2.DatesPresent = append(user2.DatesPresent, time.Now().AddDate(1, -6, 10))
	_, err = db.Model(user2).WherePK().Update()
	if err != nil {
		panic(err)
	}
	app.Run()
}
