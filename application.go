package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	db, err := gorm.Open("sqlite3", "/tmp/test.db")

	if err != nil {
		panic(err)
	}

	return &Application{
		db:     db,
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

	// [GET] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", FindNotification).
		Methods("GET")

	// [GET] /api/v1/notification
	api.HandleFunc("/notification", GetNotification).
		Methods("GET")

	// [GET] /api/v1/tag/{id}
	api.HandleFunc("/tag/{id}", FindTag).
		Methods("GET")

	// [POST] /api/v1/tag
	api.HandleFunc("/tag", PostTag).
		Methods("POST")

	// [GET] /api/v1/tag
	api.HandleFunc("/tag", GetTag).
		Methods("GET")
}
