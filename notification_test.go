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

// Tests that bad json values return a 400
func TestPostNotificationBadJSON(t *testing.T) {
	reqPayload := []byte(`{"test":}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification", server.URL), bytes.NewReader(reqPayload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 400, resp.StatusCode)
}

// Test that bulk inserting returns 202
func TestPostNotificationBulk(t *testing.T) {
	reqPayload := []byte(`[{"message": "test1"},{"message":"test2"}]`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification/bulk", server.URL), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/notification/bulk")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 202, resp.StatusCode)
}

// Test that bulk endpoint returns 400 when bad json is passed
func TestPostNotificationBulkBadJson(t *testing.T) {
	reqPayload := []byte(`[{"test":}]`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification/bulk", server.URL), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/notification/bulk")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 400, resp.StatusCode)
}

// Test that delete returns 202
func TestDeleteNotificationSuccess(t *testing.T) {
	n := &Notification{
		Message: "Testy",
	}

	App.db.Create(n)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), nil)

	if err != nil {
		t.Fatal("Request failed [DELETE] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 202, resp.StatusCode)
}

// Test that delete endpoint returns 404 if not found
func TestDeleteNotificationNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/notification/hdhdsjh", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [DELETE] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}

// Test that get returns object and 200
func TestGetNotificationSuccess(t *testing.T) {
	n := &Notification{
		Message: "Test GET",
	}

	App.db.Create(n)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	var o Notification

	json.NewDecoder(resp.Body).Decode(&o)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test GET", o.Message)
}

func TestGetNotificationNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification/not-found", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}
