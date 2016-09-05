package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server *httptest.Server
)

// Initialise Application
func init() {
	App = NewApp("5000", "test")
	App.setRoutes()

	// Make sure database migration is up to date
	Migrate()

	server = httptest.NewServer(App.router)
}

// Tests that the GET endpoint returns a 200
func TestGetReturns200(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// Tests that you can successfully save a notification via the api
func TestPostNotificationSuccess(t *testing.T) {

	reqPayload := []byte(`{"message":"A test notification", "source":"test", "read":true}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification", server.URL), bytes.NewReader(reqPayload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	var n Notification
	json.NewDecoder(resp.Body).Decode(&n)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, true, n.Read)
	assert.Equal(t, "test", n.Source)
}
