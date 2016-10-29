package router

import (
	"bytes"
	"time"
	"net/http"
	"html/template"

	"github.com/ronnrein/eps/router/middleware"
	"github.com/ronnrein/eps/utils"
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
		"Time",
		"GET",
		"/time",
		Time,
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

func Time(w http.ResponseWriter, r *http.Request) middleware.HandlerResult {
	data := struct {
		Time string
	}{
		time.Now().Format("2006-01-02T15:04:05Z"),
	}
	return middleware.HandlerResult{http.StatusOK, data, nil}
}
