package main

import (
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

// Application struct to keep our dependencies tight
type Application struct {
	db     *gorm.DB
	server *socketio.Server
	socket socketio.Socket
	router *mux.Router
	port   string
}

// Creates a new application and returns the pointer value
func NewApp(port string, dbDriver string, dbCreds string) *Application {
	db, err := gorm.Open(dbDriver, dbCreds)

	if err != nil {
		panic(err)
	}

	return &Application{
		db:     db,
		router: mux.NewRouter(),
		port:   port,
		server: ConnectSocket(),
	}
}

// Creates socket io connection
func ConnectSocket() *socketio.Server {
	sv, _ := socketio.NewServer(nil)
	return sv
}

// Starts the Application and creates a listener
func (a *Application) Run() {
	a.setRoutes()

	a.server.On("connection", func(so socketio.Socket) {
		a.socket = so
		so.Join("notify")
	})

	http.Handle("/socket.io/", a.server)
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

	// [POST] /api/v1/notification/bulk
	api.HandleFunc("/notification/bulk", PostNotifications).
		Methods("POST")

	// [GET] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", FindNotification).
		Methods("GET")

	// [DELETE] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", DeleteNotification).
		Methods("DELETE")

	// [PUT] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", PutNotification).
		Methods("PUT")

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
