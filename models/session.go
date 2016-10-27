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

type Session struct {
	gorm.Model
	Latitude 			string `gorm:"type:double"`
	Longitude 		string `gorm:"type:double"`
	DepartmentID	uint
}

type Sessions []Session

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetSessions",
			"GET",
			"/sessions",
			GetSessions,
		},
		router.Route{
			"GetSessionsSearch",
			"GET",
			"/sessions/search/{search}",
			GetSessionsSearch,
		},
		router.Route{
			"GetSession",
			"GET",
			"/sessions/{id}",
			GetSession,
		},
		router.Route{
			"AddSession",
			"POST",
			"/sessions",
			AddSession,
		},
		router.Route{
			"DeleteSession",
			"DELETE",
			"/sessions/{id}",
			DeleteSession,
		},
		router.Route{
			"UpdateSession",
			"POST",
			"/sessions/{id}",
			UpdateSession,
		},
		router.Route{
			"GetSessionMapPins",
			"GET",
			"/sessions/{id}/mappins",
			GetSessionMapPins,
		},
		router.Route{
			"GetSessionMessages",
			"GET",
			"/sessions/{id}/messages",
			GetSessionMessages,
		},
		router.Route{
			"GetSessionOperators",
			"GET",
			"/sessions/{id}/operators",
			GetSessionOperators,
		},
	)
}

func GetSessions(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	sessions := Sessions{}
	query := database.DB.Find(&sessions)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting sessions", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessions, nil}
}

func GetSession(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	session := Session{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&session, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Session not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, session, nil}
}

func AddSession(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	session := Session{}
	if err := decoder.Decode(&session); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding session", &err}
	}
	query := database.DB.Create(&session)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating session", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, session, nil}
}

func DeleteSession(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	session := Session{}
	query := database.DB.Where("ID = ?", id).Delete(&session)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting session", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Session deleted", nil}
}

func UpdateSession(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	session := Session{}
	query := database.DB.First(&session, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find session", &query.Error}
	}
	newSession := Session{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newSession); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding session", &err}
	}
	newSession.UpdatedAt = time.Now()
	query = database.DB.Model(session).Updates(&newSession)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating session", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Session updated", nil}
}

func GetSessionMapPins(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	mappins := MapPins{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("session_id = ?", id).Find(&mappins)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting map pins", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, mappins, nil}
}

func GetSessionMessages(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	messages := Messages{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("session_id = ?", id).Find(&messages)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting messages", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, messages, nil}
}

func GetSessionOperators(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	operators := Operators{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Joins("JOIN session_operators ON session_operators.operator_id = operators.id").Where("session_operators.session_id = ?", id).Find(&operators)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting operators", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, operators, nil}
}

func GetSessionsSearch(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	search := mux.Vars(r)["search"]
	sessions := Sessions{}
	query := database.DB.Joins("JOIN messages ON sessions.id = messages.session_id").Where("messages.message LIKE ?", "%"+search+"%").Or("sessions.id LIKE ?", "%"+search+"%").Group("sessions.id").Find(&sessions)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting sessions", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessions, nil}
}
