package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Endpoint to find a tag
// [GET] /api/v1/tag/{id}
func FindTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var t Tag

	if err := App.db.Where(&Tag{ExtId: params["id"]}).First(&t).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Tag not found"})
		return
	}

	jsonStr, _ := json.Marshal(&t)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to post a new tag
// [POST] /api/v1/tag
func PostTag(w http.ResponseWriter, r *http.Request) {
	var t Tag

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	App.db.Where(t).FirstOrCreate(&t)

	jsonStr, _ := json.Marshal(&t)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Endpoint to get tags
// [GET] /api/v1/tag
func GetTag(w http.ResponseWriter, r *http.Request) {
	var t []Tag

	p := GetPaginationFromRequest(r)

	App.db.Limit(p.Limit).Offset(p.Offset).Find(&t)
	jsonStr, _ := json.Marshal(&t)

	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
