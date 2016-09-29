package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"strings"
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

	Validator = validator.New()

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

// Handler to serve socket connection while added access control
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	App.server.ServeHTTP(w, r)
}

// Starts the Application and creates a listener
func (a *Application) Run() {
	a.setRoutes()

	a.server.On("connection", func(so socketio.Socket) {
		a.socket = so
		so.Join("notify")
		// We need to send notifications and tags on connect
		so.Emit("load", a.OnSocketLoad)
		so.On("notification:read", a.OnNotificationRead)
	})

	http.Handle("/socket.io/", App)
	http.Handle("/", a.router)
	log.Println("Starting Web Server on " + a.port)
	log.Fatal(http.ListenAndServe(a.port, nil))
}

// Fetches the latest notifications and tags and returns a json response
func (a *Application) OnSocketLoad() []byte {
	var t []Tag
	var n []Notification

	App.db.Find(&t)
	App.db.
		Where(&Notification{Read: false}).
		Limit(25).
		Order("created_at desc").
		Find(&n)

	p := &SocketLoadPayload{
		Notifications: n,
		Tags:          t,
	}

	return p.Serialize()
}

func (a *Application) OnNotificationRead(msg string) []byte {
	var n NotificationRead
	r := &Response{}
	err := json.NewDecoder(strings.NewReader(msg)).Decode(&n)

	if err != nil {
		r.Error = "Invalid json"
		return r.Serialize()
	}

	App.db.Model(&Notification{}).Where("ext_id IN (?)", n.Ids).Updates(map[string]interface{}{"read": true})
	r.Message = "Success"
	return r.Serialize()
}

func (a *Application) OnNotificationRefresh(msg string) []byte {
	var r NotificationRefresh
	err := json.NewDecoder(strings.NewReader(msg)).Decode(&r)

	if err != nil {
		r := &Response{Error: "Invalid json"}
		return r.Serialize()
	}

	p := GetPaginationFromSocketRequest(&r)

	var n []Notification
	var c int

	App.db.
		Model(&Notification{}).
		Where(&Notification{Read: false}).
		Count(&c)

	App.db.
		Where(&Notification{Read: false}).
		Limit(p.Limit).
		Offset(p.Offset).
		Find(&n)

	resp := &SocketLoadPayload{
		Notifications: n,
		HasNext:       (c > (r.Page * p.Limit)),
		HasPrev:       (p.Offset > 0),
	}

	return resp.Serialize()
}

// Creates the Routes
func (a *Application) setRoutes() {
	a.router.HandleFunc("/ping", PingHandler)

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

	// [POST] /api/v1/notification/{id}/read
	api.HandleFunc("/notification/{id}/read", PostReadNotification).
		Methods("POST")

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
