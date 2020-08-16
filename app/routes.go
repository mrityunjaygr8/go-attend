package app

// HandleRequests handles the API requests
func (app *App) handleRequests() {
	app.Router.HandleFunc("/api/v1/users/", app.allUsers).Methods("GET")
	app.Router.HandleFunc("/api/v1/users/", app.createUser).Methods("POST")
	app.Router.HandleFunc("/api/v1/users/{id:[0-9]+}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/api/v1/users/{id:[0-9]+}", app.deleteUser).Methods("DELETE")
}
