package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mrityunjaygr8/go-attend/users"
)

// AllUsers fetches all the users in the database
func (app *App) allUsers(w http.ResponseWriter, r *http.Request) {
	users, err := users.GetUsers(app.Db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, users)
}

// CreateUser creates a new user in the database
func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	var user users.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	if (user.Email == "") || (user.FName == "") || (user.LName == "") {
		respondWithError(w, http.StatusBadRequest, "Missing fields")
		return
	}

	err := user.CreateUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, user)
}

// GetUser fetches the specified user by ID
func (app *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := &users.User{ID: int64(id)}
	user, err = user.GetUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser deletes a user specified by ID
func (app *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := &users.User{ID: int64(id)}
	err = user.DeleteUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
