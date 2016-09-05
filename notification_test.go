package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var resp *httptest.ResponseRecorder

// Initialise Application
func init() {
	App = NewApp("5000", "test")
	App.setRoutes()

	// Make sure database migration is up to date
	Migrate()

	resp = httptest.NewRecorder()
}

// Tests that the GET endpoint returns a 200
func TestGetReturns200(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/notification", nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	App.router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}
