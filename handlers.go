package main

import (
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
