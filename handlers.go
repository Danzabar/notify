package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Method to write a REST response
func WriteResponse(w http.ResponseWriter, code int, resp RestResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp.Serialize())
}

// Endpoint to create a new notification [POST] /api/v1/notification
func PostNotification(w http.ResponseWriter, r *http.Request) {
	var n Notification

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		log.Println(err)
		return
	}

	log.Printf("%v", n)
}

// Endpoint to fetch notifications [GET] /api/v1/notification
func GetNotification(w http.ResponseWriter, r *http.Request) {

}
