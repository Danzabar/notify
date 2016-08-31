package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Endpoint to create a new notification [POST] /api/v1/notification
func PostNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var n Notification

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)

	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	log.Printf("%v", n)
}

// Endpoint to fetch notifications [GET] /api/v1/notification
func GetNotification(w http.ResponseWriter, r *http.Request) {

}
