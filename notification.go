package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Endpoint to mark a notification as read
// [POST] /api/v1/notification/{id}/read
func PostReadNotification(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	n := &Notification{}

	err := App.db.Where(&Notification{ExtId: params["id"]}).Find(n).Error

	if err != nil {
		WriteResponse(w, 404, &Response{Error: "Notification not found"})
		return
	}

	n.Read = true
	App.db.Save(n)

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
	n := NotificationRequest{}
	err := n.Deserialize(r)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	tx := App.db.Begin()

	for k := range n.Tags {
		if err := tx.Where(&n.Tags[k]).FirstOrCreate(&n.Tags[k]).Error; err != nil {
			WriteResponse(w, 422, &Response{Error: fmt.Sprintf("Unable to upsert tag - (%s)", err)})
			tx.Rollback()
			return
		}
	}

	for k := range n.Notifications {
		if err := Validator.Struct(n.Notifications[k]); err != nil {
			WriteValidationErrorResponse(w, err)
			tx.Rollback()
			return
		}

		if err := tx.Create(&n.Notifications[k]).Error; err != nil {
			WriteResponse(w, 422, &Response{Error: fmt.Sprintf("Unable to create notification[%d] - (%s)", k, err)})
			tx.Rollback()
			return
		}

		tx.Model(&n.Notifications[k]).Association("Tags").Append(n.Tags)

		if App.socket != nil {
			go func() {
				sp, _ := json.Marshal(n.Notifications[k])
				App.socket.Emit("new:notify", string(sp))
			}()
		}
	}

	tx.Commit()

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

	err = Validator.Struct(&n)

	if err != nil {
		WriteValidationErrorResponse(w, err)
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

	if err := App.db.Where("ID = ?", params["id"]).Or(&Notification{ExtId: params["id"]}).Preload("Tags").First(&n).Error; err != nil {
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
	hasRead, err := strconv.ParseBool(r.FormValue("read"))
	source := r.FormValue("source")

	s := Notification{Read: hasRead}

	if source != "" {
		s.Source = source
	}

	if err != nil {
		hasRead = false
	}

	p := GetPaginationFromRequest(r)

	App.db.Preload("Tags").
		Limit(p.Limit).
		Offset(p.Offset).
		Where(&s).
		Find(&n)

	jsonStr, _ := json.Marshal(&n)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
