package app

import (
	"encoding/json"
	"net/http"
)

// HandleRequests handles the API requests
func (app *App) handleRequests() {
	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(jsonResponse{Message: "Hello"})
	})

	auth := app.Router.PathPrefix("/api/v1/auth").Subrouter()
	auth.HandleFunc("/login", app.login).Methods("POST")
	auth.HandleFunc("/register", app.createUser).Methods("POST")

	api := app.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(jwtVerify)
	api.HandleFunc("/users", app.allUsers).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", app.getUser).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", app.deleteUser).Methods("DELETE")
	api.HandleFunc("/users/{id:[0-9]+}", app.updateUser).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}/password", app.updatePassword).Methods("PUT")
}
