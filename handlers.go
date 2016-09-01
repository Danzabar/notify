package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Method to write a REST response
func WriteResponse(w http.ResponseWriter, code int, resp RestResponse) {
	WriteResponseHeader(w, code)
	w.Write(resp.Serialize())
}

// Writes the headers for the response
func WriteResponseHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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

	App.db.Create(&n)

	jsonStr, _ := json.Marshal(&n)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to fetch notifications [GET] /api/v1/notification
func GetNotification(w http.ResponseWriter, r *http.Request) {
	var n []Notification

	App.db.Find(&n)
	jsonStr, err := json.Marshal(&n)

	if err != nil {
		WriteResponse(w, 500, &Response{Error: "Unable to marshal models"})
		return
	}

	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
