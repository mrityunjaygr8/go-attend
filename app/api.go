package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mrityunjaygr8/go-attend/users"
	"golang.org/x/crypto/bcrypt"
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

	if (user.Email == "") || (user.FName == "") || (user.LName == "") || (user.Role == "") || (user.Password == "") {
		respondWithError(w, http.StatusBadRequest, "Missing fields")
		return
	}

	if !inArray(user.Role, []string{"base", "admin"}) {
		respondWithError(w, http.StatusBadRequest, "The 'Role' field must be one of: base, admin")
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

func (app *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	user := &users.User{ID: int64(id)}
	user, err = user.GetUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	var newUser users.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	if (newUser.Email != user.Email) || (newUser.FName != user.FName) || (newUser.LName != user.LName) || (newUser.Password != user.Password) {
		respondWithError(w, http.StatusBadRequest, "Cannot edit the email, first name, last name and password.")
		return
	}

	if !inArray(newUser.Role, []string{"base", "admin"}) {
		respondWithError(w, http.StatusBadRequest, "The 'Role' field must be one of: base, admin")
		return
	}

	newUser.ID = user.ID

	err = newUser.UpdateUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, newUser)
}

func (app *App) updatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	user := &users.User{ID: int64(id)}
	user, err = user.GetUser(app.Db)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	var password users.UserPassChange
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&password); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	user.Password = password.Password

	err = user.UpdatePassword(app.Db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
	user := &users.UserLoginStruct{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	resp, err := app.findOne(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (app *App) findOne(email, password string) (map[string]interface{}, error) {
	user, err := users.GetForLogin(app.Db, email)
	if err != nil {
		return nil, err
	}

	// return user, nil

	expiresAt := time.Now().Add(time.Minute * 10000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("Invalid User Credentials")
	}

	tk := users.JWTToken{
		Email: user.Email,
		FName: user.FName,
		LName: user.LName,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	var response = map[string]interface{}{}
	response["token"] = tokenString
	response["user"] = user

	return response, nil
}
