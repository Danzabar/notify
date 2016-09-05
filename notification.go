package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Endpoint to bulk insert notifications
// [POST] /api/v1/notification/bulk
func PostNotifications(w http.ResponseWriter, r *http.Request) {
	var n []Notification

	err := json.NewDecoder(r.Body).Decode(&n)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	go func() {
		for _, r := range n {
			App.db.Create(&r)
		}
	}()

	WriteResponseHeader(w, 202)
}

// Endpoint to delete notifications
// [DELETE] /api/v1/notification/{id}
func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var n Notification

	err := App.db.Where(&Notification{ExtId: params["id"]}).Find(&n).Error

	if err != nil {
		WriteResponse(w, 404, &Response{Error: "Notification not found"})
		return
	}

	App.db.Delete(&n)
	WriteResponseHeader(w, 202)
}

// Endpoint to create a new notification
// [POST] /api/v1/notification
func PostNotification(w http.ResponseWriter, r *http.Request) {
	var n Notification

	err := json.NewDecoder(r.Body).Decode(&n)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	if err := App.db.Create(&n).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save notification"})
		return
	}

	jsonStr, _ := json.Marshal(&n)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to update a notification
// [PUT] /api/v1/notification/{id}
func PutNotification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var n Notification

	if err := App.db.Where("ID = ?", params["id"]).Or(&Notification{ExtId: params["id"]}).First(&n).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Notification not found"})
		return
	}

	var u Notification
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	if err := App.db.Model(&n).Updates(u).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save notification"})
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
	jsonStr, _ := json.Marshal(&n)

	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
