package main

import (
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

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

func WriteValidationErrorResponse(w http.ResponseWriter, err error) {
	v := &ValidationResponse{
		Errors: make(map[string]string),
	}

	for _, e := range err.(validator.ValidationErrors) {
		v.Errors[e.Field()] = "This field is invalid"
	}

	WriteResponse(w, 400, v)
}

// Creates pagination options from given request
func GetPaginationFromRequest(r *http.Request) *Pagination {
	p := &Pagination{}

	limit, err := strconv.ParseInt(r.FormValue("pageSize"), 10, 8)

	if err != nil {
		p.Limit = 50
	} else {
		p.Limit = int(limit)
	}

	offset, err := strconv.ParseInt(r.FormValue("page"), 10, 8)

	if err != nil {
		offset = 1
	}

	p.Offset = (int(offset) - 1) * p.Limit
	return p
}
