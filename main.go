package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrityunjaygr8/go-attend/app"
)

func main() {
	godotenv.Load()

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

	app.Run()
}
