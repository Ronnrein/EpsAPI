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

type MapPin struct {
	gorm.Model
	Latitude 	string `gorm:"type:double"`
	Longitude string `gorm:"type:double"`
	Text			string
	SessionID	uint
}

type MapPins []MapPin

func init() {
	router.Routes = append(
		router.Routes,
		router.Route{
			"GetMapPins",
			"GET",
			"/mappins",
			GetMapPins,
		},
		router.Route{
			"GetMapPin",
			"GET",
			"/mappins/{id}",
			GetMapPin,
		},
		router.Route{
			"AddMapPin",
			"POST",
			"/mappins",
			AddMapPin,
		},
		router.Route{
			"DeleteMapPin",
			"DELETE",
			"/mappins/{id}",
			DeleteMapPin,
		},
		router.Route{
			"UpdateMapPin",
			"POST",
			"/mappins/{id}",
			UpdateMapPin,
		},
	)
}

func GetMapPins(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	mappins := MapPins{}
	query := database.DB.Find(&mappins)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting map pins", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, mappins, nil}
}

func GetMapPin(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	idStr := mux.Vars(r)["id"]
	mappin := MapPin{}
	var query *gorm.DB
	if id, err := strconv.Atoi(idStr); err == nil {
		query = database.DB.First(&mappin, id)
	} else {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Map pin not found", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, mappin, nil}
}

func AddMapPin(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	decoder := json.NewDecoder(r.Body)
	mappin := MapPin{}
	if err := decoder.Decode(&mappin); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding map pin", &err}
	}
	query := database.DB.Create(&mappin)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error creating map pin", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, mappin, nil}
}

func DeleteMapPin(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	mappin := MapPin{}
	query := database.DB.Where("ID = ?", id).Delete(&mappin)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error deleting map pin", &query.Error}

	}
	return middleware.HandlerResult{http.StatusOK, "Map pin deleted", nil}
}

func UpdateMapPin(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Bad ID supplied", &err}
	}
	mappin := MapPin{}
	query := database.DB.First(&mappin, id)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusNotFound, "Could not find map pin", &query.Error}
	}
	newMapPin := MapPin{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&newMapPin); err != nil {
		return middleware.HandlerResult{http.StatusBadRequest, "Error decoding map pin", &err}
	}
	newMapPin.UpdatedAt = time.Now()
	query = database.DB.Model(mappin).Updates(&newMapPin)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error updating map pin", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, "Map pin updated", nil}
}
