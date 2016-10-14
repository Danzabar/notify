package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// [POST] /api/v1/alert-group
func PostAlertGroup(w http.ResponseWriter, r *http.Request) {
	a := &AlertGroupRequest{}
	if err := a.Deserialize(r); err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	err := Validator.Struct(&a.Group)

	if err != nil {
		WriteValidationErrorResponse(w, err)
		return
	}

	tx := App.db.Begin()

	if err := tx.Create(&a.Group).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save entity"})
		tx.Rollback()
		return
	}

	for k := range a.Recipients {
		tx.Where(&a.Recipients[k]).FirstOrCreate(&a.Recipients[k])
		tx.Model(&a.Group).Association("Recipients").Append(&a.Recipients[k])
	}

	for k := range a.Tags {
		tx.Where(&a.Tags[k]).FirstOrCreate(&a.Tags[k])
		tx.Model(&a.Group).Association("Tags").Append(&a.Tags[k])
	}

	tx.Commit()

	jsonStr, _ := json.Marshal(&a.Group)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// [PUT] /api/v1/alert-group/{id}
func PutAlertGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var a AlertGroup
	var u AlertGroupRequest

	if err := App.db.Where("ext_id = ?", params["id"]).First(&a).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Alert group not found"})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	err = Validator.Struct(&u.Group)

	if err != nil {
		WriteValidationErrorResponse(w, err)
		return
	}

	if err := App.db.Model(&a).Updates(u.Group).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save entity"})
		return
	}

	jsonStr, _ := json.Marshal(&a)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// [GET] /api/v1/alert-group
func GetAlertGroup(w http.ResponseWriter, r *http.Request) {
	var a []AlertGroup

	p := GetPaginationFromRequest(r)
	App.db.Limit(p.Limit).Offset(p.Offset).Find(&a)

	jsonStr, _ := json.Marshal(&a)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// [DELETE] /api/v1/alert-group/{id}
func DeleteAlertGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var a AlertGroup

	if err := App.db.Where("ext_id = ?", params["id"]).First(&a).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Alert group not found"})
		return
	}

	App.db.Delete(&a)
	WriteResponseHeader(w, 202)
}
