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

type Operator struct {
	gorm.Model
	Name					string
	Username			string
	Password			string
	DepartmentID	uint
}

type Operators []Operator

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetOperators",
			"GET",
			"/operators",
			GetOperators,
		},
		router.Route{
			"GetOperator",
			"GET",
			"/operators/{id}",
			GetOperator,
		},
		router.Route{
			"AddOperator",
			"POST",
			"/operators",
			AddOperator,
		},
		router.Route{
			"DeleteOperator",
			"DELETE",
			"/operators/{id}",
			DeleteOperator,
		},
		router.Route{
			"UpdateOperator",
			"POST",
			"/operators/{id}",
			UpdateOperator,
		},
		router.Route{
			"GetOperatorMessages",
			"GET",
			"/operators/{id}/messages",
			GetOperatorMessages,
		},
		router.Route{
			"GetOperatorSessions",
			"GET",
			"/operators/{id}/sessions",
			GetOperatorSessions,
		},
	)
}

func GetOperators(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	operators := Operators{}
	query := database.DB.Find(&operators)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting operators", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, operators, nil}
}

func GetOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	operator := Operator{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&operator, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Operator not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, operator, nil}
}

func AddOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	operator := Operator{}
	if err := decoder.Decode(&operator); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding operator", &err}
	}
	query := database.DB.Create(&operator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating operator", &query.Error}
	}
	url := fmt.Sprintf("http://%s:%d/operators/%d", utils.Config.Host, utils.Config.Port, operator.ID)
	return middleware.HandlerResult{http.StatusOK, url, nil}
}

func DeleteOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	operator := Operator{}
	query := database.DB.Where("ID = ?", id).Delete(&operator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting operator", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Operator deleted", nil}
}

func UpdateOperator(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	operator := Operator{}
	query := database.DB.First(&operator, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find operator", &query.Error}
	}
	newOperator := Operator{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newOperator); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding operator", &err}
	}
	newOperator.UpdatedAt = time.Now()
	query = database.DB.Model(operator).Updates(&newOperator)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating operator", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Operator updated", nil}
}

func GetOperatorMessages(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	messages := Messages{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("operator_id = ?", id).Find(&messages)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting messages", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, messages, nil}
}

func GetOperatorSessions(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	sessions := Sessions{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Joins("JOIN session_operators ON session_operators.session_id = sessions.id").Where("session_operators.operator_id = ?", id).Find(&sessions);
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting sessions", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessions, nil}
}
