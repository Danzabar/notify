package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		log.Println(err)
		return
	}

	if err := App.db.Where(t).FirstOrCreate(&t).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save tag"})
		log.Println(err)
		return
	}

	jsonStr, _ := json.Marshal(&t)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

func GetTag(w http.ResponseWriter, r *http.Request) {
	var t []Tag

	App.db.Find(&t)
	jsonStr, err := json.Marshal(&t)

	if err != nil {
		WriteResponse(w, 500, &Response{Error: "Unable to marshal models"})
		return
	}

	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
