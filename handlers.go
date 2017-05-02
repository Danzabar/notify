package main

import (
	"encoding/base64"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
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

// Used to chain handlers with middleware
func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// Handler for Basic Auth
func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not Authorized", 401)
			return
		}

		b, _ := base64.StdEncoding.DecodeString(s[1])
		pair := strings.SplitN(string(b), ":", 2)

		if len(pair) != 2 {
			http.Error(w, "Not Authorized", 401)
			return
		}

		if pair[0] != App.user || pair[1] != App.pass {
			http.Error(w, "Not Authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}

// Creates pagination options from given request
func GetPaginationFromRequest(r *http.Request) *Pagination {
	var l int
	var o int
	limit, err := strconv.ParseInt(r.FormValue("pageSize"), 10, 8)

	if err != nil {
		l = 50
	} else {
		l = int(limit)
	}

	offset, err := strconv.ParseInt(r.FormValue("page"), 10, 8)

	if err != nil {
		o = 1
	} else {
		o = int(offset)
	}

	return createPagination(l, o)
}

func GetPaginationFromSocketRequest(r NotificationRefresh) *Pagination {
	if r.Page == 0 {
		r.Page = 1
	}

	if r.PageSize == 0 {
		r.PageSize = 50
	}

	return createPagination(r.PageSize, r.Page)
}

func createPagination(limit int, offset int) *Pagination {
	p := &Pagination{}

	p.Limit = limit
	p.Offset = (offset - 1) * p.Limit

	return p
}
