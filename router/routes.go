package router

import (
	"bytes"
	"net/http"
	"html/template"

	"github.com/ronnrein/eps/router/middleware"
	"github.com/ronnrein/eps/utils"
	"github.com/ronnrein/eps/database"
	"github.com/gorilla/mux"
	"strconv"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc middleware.ResultHandler
}

var HtmlRoutes []Route

var Routes = []Route{
	Route{
		"Index", // Name
		"GET",   // Method
		"/",     // Pattern
		Index,   // Handler
	},
	Route{
		"Log",
		"GET",
		"/log",
		Log,
	},
	Route{
		"LogLimit",
		"GET",
		"/log/{id}",
		Log,
	},
}

func Index(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	t, _ := template.ParseFiles("static/index.html")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	data := struct {
		Routes []Route
		Config utils.Conf
	}{
		HtmlRoutes,
		utils.Config,
	}
	var doc bytes.Buffer
	if err := t.Execute(&doc, data); err != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error generating index site from template", &err}
	}

	return middleware.HandlerResult{http.StatusOK, doc.String(), nil}
}

func Log(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	id := mux.Vars(r)["id"]
	limit, err := strconv.Atoi(id)
	if err != nil {
		limit = 100
	}
	logs := []middleware.LogEntry{}
	query := database.DB.Limit(limit).Find(&logs)
	if query.Error != nil {
		return middleware.HandlerResult{http.StatusInternalServerError, "Error getting logs", &query.Error}
	}
	return middleware.HandlerResult{http.StatusOK, logs, nil}
}
