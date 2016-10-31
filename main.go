package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ronnrein/eps/database"
	"github.com/ronnrein/eps/models"
	"github.com/ronnrein/eps/router"
	"github.com/ronnrein/eps/router/middleware"
	"github.com/ronnrein/eps/utils"
)

func init() {
	database.DB.AutoMigrate(&models.Department{}, &models.MapPin{}, &models.Message{}, &models.Operator{}, &models.Session{}, &models.SessionOperator{}, &middleware.LogEntry{})
}

func main() {
	router.RemoveHandlers()
	router := router.NewRouter()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprint(":", utils.Config.SSLPort), utils.Config.SSLCert, utils.Config.SSLKey, router))
}
