package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func FindTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var t Tag

	if err := App.db.Where("ID = ?", params["id"]).First(&t).Error; err != nil {
		WriteResponse(w, 404, &Response{Error: "Tag not found"})
		return
	}

	jsonStr, _ := json.Marshal(&t)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

func PostTag(w http.ResponseWriter, r *http.Request) {
	var t Tag

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	if err := App.db.Where(t).FirstOrCreate(&t).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save tag"})
		return
	}

	jsonStr, _ := json.Marshal(&t)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

func GetTag(w http.ResponseWriter, r *http.Request) {
	var t []Tag

	App.db.Find(&t)
	jsonStr, _ := json.Marshal(&t)

	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
