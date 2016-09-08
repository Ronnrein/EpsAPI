package models

import (
	"fmt"
	"strconv"
	"time"
	"net/http"
	"encoding/json"

	"github.com/ronnrein/eps/database"
	"github.com/ronnrein/eps/router"
	"github.com/ronnrein/eps/utils"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/ronnrein/eps/router/middleware"
)

type User struct {
	gorm.Model
	DeviceID string
}

type Users []User

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetUsers",
			"GET",
			"/users",
			GetUsers,
		},
		router.Route{
			"GetUser",
			"GET",
			"/users/{id}",
			GetUser,
		},
		router.Route{
			"AddUser",
			"POST",
			"/users",
			AddUser,
		},
		router.Route{
			"UpdateUser",
			"POST",
			"/users/{id}",
			UpdateUser,
		},
		router.Route{
			"GetUserSessions",
			"GET",
			"/users/{id}/sessions",
			GetUserSessions,
		},
	)
}

func GetUsers(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	users := Users{}
	query := database.DB.Find(&users)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting users", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, users, nil}
}

func GetUser(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	user := User{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&user, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "User not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, user, nil}
}

func AddUser(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	if err := decoder.Decode(&user); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding user", &err}
	}
	query := database.DB.Create(&user)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating user", &query.Error}
	}
	url := fmt.Sprintf("http://%s:%d/users/%d", utils.Config.Host, utils.Config.Port, user.ID)
	return middleware.HandlerResult{http.StatusOK, url, nil}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	user := User{}
	query := database.DB.Where("id = ?", id).Delete(&user)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting user", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "User deleted", nil}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	user := User{}
	query := database.DB.First(&user, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find user", &query.Error}
	}
	newUser := User{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newUser); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding user", &err}
	}
	newUser.UpdatedAt = time.Now()
	query = database.DB.Model(user).Updates(&newUser)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating user", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "User updated", nil}
}

func GetUserSessions(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	sessions := Sessions{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("user_id = ?", id).Find(&sessions)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting sessions", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessions, nil}
}
