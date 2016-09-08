package router

import (
	"github.com/gorilla/mux"

	"github.com/ronnrein/eps/router/middleware"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range Routes {
		handler := useHandlers(route.HandlerFunc, middleware.Logger, middleware.SetHeaders)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
