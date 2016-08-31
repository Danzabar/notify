package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

// Application struct to keep our dependencies tight
type Application struct {
	db     *gorm.DB
	router *mux.Router
	port   string
}

// Creates a new application and returns the pointer value
func NewApp(port string) *Application {
	return &Application{
		router: mux.NewRouter(),
		port:   port,
	}
}

// Starts the Application and creates a listener
func (a *Application) Run() {
	a.setRoutes()

	http.Handle("/", a.router)
	log.Println("Starting Web Server on " + a.port)
	log.Fatal(http.ListenAndServe(a.port, nil))
}

// Creates the Routes
func (a *Application) setRoutes() {
	// API specific routes
	api := a.router.PathPrefix("/api/v1").Subrouter()

	// [POST] /api/v1/notification
	api.HandleFunc("/notification", PostNotification).
		Methods("POST")

	// [GET] /api/v1/notification
	api.HandleFunc("/notification", GetNotification).
		Methods("GET")
}
