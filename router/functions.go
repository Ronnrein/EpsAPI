package router

import (
	"strings"
	"net/http"
)

func RemoveHandlers() {
	var newRoutes []Route
	for _, route := range Routes {
		pattern := strings.Replace(route.Pattern, "{id}", "1", -1)
		newRoute := Route{Name: route.Name, Method: route.Method, Pattern: pattern}
		newRoutes = append(newRoutes, newRoute)
	}
	HtmlRoutes = newRoutes
}

func useHandlers(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	var res http.Handler = h
	for _, m := range middleware {
		res = m(res)
	}

	return res
}
