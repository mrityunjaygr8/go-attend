package app

import (
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/mrityunjaygr8/go-attend/users"
)

// App is the database interface and the App
type App struct {
	Db     *pg.DB
	Port   string
	Router *mux.Router
}

type jsonResponse struct {
	Message string `json:"message"`
}

// CreateSchema creates the DB schema for the models
func (app *App) CreateSchema() error {
	models := []interface{}{
		(*users.User)(nil),
	}

	for _, model := range models {
		err := app.Db.Model(model).CreateTable(&orm.CreateTableOptions{})

		if err != nil {
			return err
		}

	}
	return nil
}

// Setup sets up the database
func (app *App) Setup(config string) error {
	opt, err := pg.ParseURL(config)
	if err != nil {
		return err
	}

	app.Db = pg.Connect(opt)
	// err = app.CreateSchema()
	// if err != nil {
	// 	return err
	// }
	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router.Use(loggingMiddleware)
	app.handleRequests()
	return nil
}

// Close closes the database connection
func (app *App) Close() {
	app.Db.Close()
}

// Run the server/router
func (app *App) Run() {
	log.Fatal(http.ListenAndServe(app.Port, app.Router))
}
