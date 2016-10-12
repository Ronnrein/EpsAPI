package models

import (
	"strconv"
	"time"
	"net/http"
	"encoding/json"

	"github.com/ronnrein/eps/database"
	"github.com/ronnrein/eps/router"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/ronnrein/eps/router/middleware"
)

type Message struct {
	gorm.Model
	Message			string
	OperatorID	uint
	SessionID		uint
}

type Messages []Message

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetMessages",
			"GET",
			"/messages",
			GetMessages,
		},
		router.Route{
			"GetMessage",
			"GET",
			"/messages/{id}",
			GetMessage,
		},
		router.Route{
			"AddMessage",
			"POST",
			"/messages",
			AddMessage,
		},
		router.Route{
			"DeleteMessage",
			"DELETE",
			"/messages/{id}",
			DeleteMessage,
		},
		router.Route{
			"UpdateMessage",
			"POST",
			"/messages/{id}",
			UpdateMessage,
		},
	)
}

func GetMessages(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	messages := Messages{}
	query := database.DB.Find(&messages)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting messages", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, messages, nil}
}

func GetMessage(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	message := Message{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&message, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Message not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, message, nil}
}

func AddMessage(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	message := Message{}
	if err := decoder.Decode(&message); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding message", &err}
	}
	query := database.DB.Create(&message)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating message", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, message, nil}
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	message := Message{}
	query := database.DB.Where("ID = ?", id).Delete(&message)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting message", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Message deleted", nil}
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	message := Message{}
	query := database.DB.First(&message, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find message", &query.Error}
	}
	newMessage := Message{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newMessage); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding message", &err}
	}
	newMessage.UpdatedAt = time.Now()
	query = database.DB.Model(message).Updates(&newMessage)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating message", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Message updated", nil}
}
