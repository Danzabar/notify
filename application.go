package main

import (
	"encoding/json"
	"github.com/Danzabar/WatchDog/plugins"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
	"gopkg.in/go-playground/validator.v9"
	translations "gopkg.in/go-playground/validator.v9/translations/en"
	"net/http"
	"os"
	"strings"
)

// Application struct to keep our dependencies tight
type Application struct {
	db     *gorm.DB
	server *socketio.Server
	socket socketio.Socket
	router *mux.Router
	mg     Alerter
	pb     Alerter
	log    *logging.Logger
	trans  ut.Translator
	test   bool
	port   string
	user   string
	pass   string
}

type Alerter interface {
	SendNotification(a AlertGroup, t string) bool
}

// Mock for alerting
type MockAlerter struct {
	Pass bool
}

func (m *MockAlerter) SendNotification(a AlertGroup, t string) bool {
	return m.Pass
}

// Creates a new application and returns the pointer value
func NewApp(port string, dbDriver string, dbCreds string, user string, pass string) *Application {
	db, err := gorm.Open(dbDriver, dbCreds)

	if err != nil {
		panic(err)
	}

	SetLogging()
	t := CreateValidator()

	return &Application{
		db:     db,
		router: mux.NewRouter(),
		port:   port,
		server: ConnectSocket(),
		user:   user,
		log:    logging.MustGetLogger("notify"),
		pass:   pass,
		mg:     NewMailClient(),
		pb:     NewPushBullet(),
		trans:  t,
	}
}

func CreateValidator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	Validator = validator.New()

	translations.RegisterDefaultTranslations(Validator, trans)
	return trans
}

// Sets up settings for logs
func SetLogging() {
	f := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	b := logging.NewLogBackend(os.Stderr, "", 0)

	bf := logging.NewBackendFormatter(b, f)
	logging.SetBackend(bf)
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
		a.log.Debug("User connects to socket")
		App.socket = so
	})

	a.server.On("notification:read", a.OnNotificationRead)
	a.server.On("notification:refresh", a.OnNotificationRefresh)

	http.Handle("/socket.io/", App)
	http.Handle("/", a.router)
	a.log.Debugf("Starting Web Server on " + a.port)
	a.log.Fatal(http.ListenAndServe(a.port, nil))
}

func (a *Application) OnNotificationRead(msg string) string {
	var n NotificationRead
	r := &Response{}
	err := json.NewDecoder(strings.NewReader(msg)).Decode(&n)

	if err != nil {
		r.Error = "Invalid json"
		return string(r.Serialize())
	}

	App.db.Model(&Notification{}).Where("ext_id IN (?)", n.Ids).Updates(map[string]interface{}{"read": true})
	r.Message = "Success"
	return string(r.Serialize())
}

func (a *Application) OnNotificationRefresh(msg string) string {
	var r NotificationRefresh
	err := json.NewDecoder(strings.NewReader(msg)).Decode(&r)

	if err != nil {
		r := &Response{Error: "Invalid json"}
		return string(r.Serialize())
	}

	p := GetPaginationFromSocketRequest(r)
	c, n := SearchNotifications(r, p)

	resp := &SocketLoadPayload{
		Notifications: n,
		HasNext:       (c > (r.Page * p.Limit)),
		HasPrev:       (p.Offset > 0),
	}

	return string(resp.Serialize())
}

// Creates the Routes
func (a *Application) setRoutes() {

	// API specific routes
	api := a.router.PathPrefix("/api/v1").Subrouter()

	a.router.HandleFunc("/health", Use(plugins.HealthCheckHandler, BasicAuth))

	// [POST] /api/v1/alert-group
	api.HandleFunc("/alert-group", PostAlertGroup).
		Methods("POST")

	// [GET] /api/v1/alert-group
	api.HandleFunc("/alert-group", Use(GetAlertGroup, BasicAuth)).
		Methods("GET")

	// [PUT] /api/v1/alert-group/{id}
	api.HandleFunc("/alert-group/{id}", PutAlertGroup).
		Methods("PUT")

	// [DELETE] /api/v1/alert-group/{id}
	api.HandleFunc("/alert-group/{id}", DeleteAlertGroup).
		Methods("DELETE")

	// [POST] /api/v1/notification
	api.HandleFunc("/notification", PostNotification).
		Methods("POST")

	// [GET] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", Use(FindNotification, BasicAuth)).
		Methods("GET")

	// [DELETE] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", Use(DeleteNotification, BasicAuth)).
		Methods("DELETE")

	// [PUT] /api/v1/notification/{id}
	api.HandleFunc("/notification/{id}", PutNotification).
		Methods("PUT")

	// [POST] /api/v1/notification/{id}/read
	api.HandleFunc("/notification/{id}/read", PostReadNotification).
		Methods("POST")

	// [GET] /api/v1/notification
	api.HandleFunc("/notification", Use(GetNotification, BasicAuth)).
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
