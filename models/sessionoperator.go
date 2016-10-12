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

type SessionOperator struct {
	gorm.Model
	SessionID 	uint
	OperatorID	uint
}

type SessionOperators []SessionOperator

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"AddSessionOperator",
			"POST",
			"/sessionoperators",
			AddSessionOperator,
		},
		router.Route{
			"DeleteSessionOperator",
			"DELETE",
			"/sessionoperators/{id}",
			DeleteSessionOperator,
		},
		router.Route{
			"UpdateSessionOperator",
			"POST",
			"/sessionoperators/{id}",
			UpdateSessionOperator,
		},
	)
}

func AddSessionOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	sessionoperator := SessionOperator{}
	if err := decoder.Decode(&sessionoperator); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding session operator", &err}
	}
	query := database.DB.Create(&sessionoperator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating session operator", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessionoperator, nil}
}

func DeleteSessionOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	sessionoperator := SessionOperator{}
	query := database.DB.Where("ID = ?", id).Delete(&sessionoperator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting session operator", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Session operator deleted", nil}
}

func UpdateSessionOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	sessionoperator := SessionOperator{}
	query := database.DB.First(&sessionoperator, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find session operator", &query.Error}
	}
	newSessionOperator := SessionOperator{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newSessionOperator); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding session operator", &err}
	}
	newSessionOperator.UpdatedAt = time.Now()
	query = database.DB.Model(sessionoperator).Updates(&newSessionOperator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating session operator", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Session operator updated", nil}
}
