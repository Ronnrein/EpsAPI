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

type Department struct {
	gorm.Model
	Name			string
	Sessions 	Sessions
	Operators	Operators
}

type Departments []Department

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetDepartments",
			"GET",
			"/departments",
			GetDepartments,
		},
		router.Route{
			"GetDepartment",
			"GET",
			"/departments/{id}",
			GetDepartment,
		},
		router.Route{
			"AddDepartment",
			"POST",
			"/departments",
			AddDepartment,
		},
		router.Route{
			"DeleteDepartment",
			"DELETE",
			"/departments/{id}",
			DeleteDepartment,
		},
		router.Route{
			"UpdateDepartment",
			"POST",
			"/departments/{id}",
			UpdateDepartment,
		},
		router.Route{
			"GetDepartmentSessions",
			"GET",
			"/departments/{id}/sessions",
			GetDepartmentSessions,
		},
		router.Route{
			"GetDepartmentOperators",
			"GET",
			"/departments/{id}/operators",
			GetDepartmentOperators,
		},
	)
}

func GetDepartments(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	departments := Departments{}
	query := database.DB.Find(&departments)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting departments", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, departments, nil}
}

func GetDepartment(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	department := Department{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&department, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Department not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, department, nil}
}

func AddDepartment(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	department := Department{}
	if err := decoder.Decode(&department); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding department", &err}
	}
	query := database.DB.Create(&department)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating department", &query.Error}
	}
	url := fmt.Sprintf("http://%s:%d/departments/%d", utils.Config.Host, utils.Config.Port, department.ID)
	return middleware.HandlerResult{http.StatusOK, url, nil}
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	department := Department{}
	query := database.DB.Where("ID = ?", id).Delete(&department)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting department", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Department deleted", nil}
}

func UpdateDepartment(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	department := Department{}
	query := database.DB.First(&department, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find department", &query.Error}
	}
	newDepartment := Department{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newDepartment); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding department", &err}
	}
	newDepartment.UpdatedAt = time.Now()
	query = database.DB.Model(department).Updates(&newDepartment)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating department", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Department updated", nil}
}

func GetDepartmentSessions(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	sessions := Sessions{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("department_id = ?", id).Find(&sessions)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting sessions", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, sessions, nil}
}

func GetDepartmentOperators(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	operators := Operators{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.Where("department_id = ?", id).Find(&operators)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting operators", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, operators, nil}
}
