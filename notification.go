package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func PostNotifications(w http.ResponseWriter, r *http.Request) {
	var n []Notification

	err := json.NewDecoder(r.Body).Decode(&n)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		log.Println(err)
		return
	}

	for _, r := range n {
		App.db.Create(&r)
	}

	WriteResponseHeader(w, 202)
}

// Endpoint to create a new notification
// [POST] /api/v1/notification
func PostNotification(w http.ResponseWriter, r *http.Request) {
	var n Notification

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		log.Println(err)
		return
	}

	if err := App.db.Create(&n).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save notification"})
		log.Println(err)
		return
	}

	jsonStr, _ := json.Marshal(&n)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to find a notification
// [GET] /api/v1/notification/{id}
func FindNotification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var n Notification

	if err := App.db.Where("ID = ?", params["id"]).Or(&Notification{ExtId: params["id"]}).First(&n).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Notification not found"})
		return
	}

	jsonStr, _ := json.Marshal(&n)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to fetch notifications
// [GET] /api/v1/notification
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
