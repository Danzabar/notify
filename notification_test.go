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
	App = NewApp("5000", "sqlite3", "/tmp/test.db", "test", "test")
	App.setRoutes()

	// Make sure database migration is up to date
	Migrate()

	server = httptest.NewServer(App.router)
}

// Tests that the GET endpoint returns a 200
func TestGetReturns200(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification", server.URL), nil)
	req.SetBasicAuth("test", "test")

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetSeperatesReadMessages(t *testing.T) {
	App.db.Create(&Notification{Message: "Test1", Read: false})
	App.db.Create(&Notification{Message: "Test2", Read: true})

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification?read=true", server.URL), nil)
	req.SetBasicAuth("test", "test")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

// Tests that you can successfully save a notification via the api
func TestPostNotificationSuccess(t *testing.T) {
	r := &NotificationRequest{}
	r.Notifications = append(r.Notifications, Notification{Message: "a test notification", Source: "test", Read: true})
	reqPayload, _ := json.Marshal(r)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification", server.URL), bytes.NewReader(reqPayload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	var n NotificationRequest
	json.NewDecoder(resp.Body).Decode(&n)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, true, n.Notifications[0].Read)
	assert.Equal(t, "test", n.Notifications[0].Source)
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

func TestPostNotificationValidation(t *testing.T) {
	r := &NotificationRequest{}
	r.Notifications = append(r.Notifications, Notification{Message: "a test notification"})
	reqPayload, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification", server.URL), bytes.NewReader(reqPayload))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

// Test that delete returns 202
func TestDeleteNotificationSuccess(t *testing.T) {
	n := &Notification{
		Message: "Testy",
	}

	App.db.Create(n)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), nil)
	req.SetBasicAuth("test", "test")

	if err != nil {
		t.Fatal("Request failed [DELETE] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 202, resp.StatusCode)
}

// Test that delete endpoint returns 404 if not found
func TestDeleteNotificationNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/notification/hdhdsjh", server.URL), nil)
	req.SetBasicAuth("test", "test")

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
	req.SetBasicAuth("test", "test")

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	var o Notification

	json.NewDecoder(resp.Body).Decode(&o)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test GET", o.Message)
}

// Test that get returns 404 if not found
func TestGetNotificationNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification/not-found", server.URL), nil)
	req.SetBasicAuth("test", "test")
	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}

// Test that put updates object and returns 200
func TestPutNotificationSuccess(t *testing.T) {
	n := &Notification{Message: "PUTTEST", Source: "Test"}
	reqPayload := []byte(`{"message":"a new message"}`)

	App.db.Create(n)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [PUT] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	var o Notification

	json.NewDecoder(resp.Body).Decode(&o)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "a new message", o.Message)
}

func TestPutNotificationValidation(t *testing.T) {
	n := &Notification{Source: "test"}
	reqPayload := []byte(`{"message":""}`)

	App.db.Create(n)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), bytes.NewReader(reqPayload))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

// Test that a put on an unknown notification returns a 404
func TestPutNotificationNotFound(t *testing.T) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/notification/not-found", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [PUT] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestPutNotificationBadJson(t *testing.T) {
	n := &Notification{Message: "PUTJSONTEST"}
	reqPayload := []byte(`{"test":}`)

	App.db.Create(n)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/notification/%s", server.URL, n.ExtId), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [PUT] /api/v1/notification")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetReturnsPaginatedList(t *testing.T) {
	App.db.Delete(&Notification{})

	for i := 0; i < 10; i++ {
		n := &Notification{Message: fmt.Sprintf("test_%d", i)}
		App.db.Create(n)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification?pageSize=1&page=5", server.URL), nil)
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var output []Notification
	json.NewDecoder(resp.Body).Decode(&output)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, len(output))
	assert.Equal(t, "Test_5", output[0].Message)
}

func TestPostReadNotificationSuccess(t *testing.T) {
	n := &Notification{Message: "Test", Source: "Test"}
	App.db.Create(n)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification/%s/read", server.URL, n.ExtId), nil)
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var not Notification
	App.db.Where(&Notification{ExtId: n.ExtId}).First(&not)

	assert.Equal(t, 202, resp.StatusCode)
	assert.Equal(t, true, not.Read)
}

func TestPostReadFail(t *testing.T) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/notification/test-fhf/read", server.URL), nil)
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetNotificationBySource(t *testing.T) {
	App.db.Create(&Notification{Message: "Test1", Source: "source1"})
	App.db.Create(&Notification{Message: "Test2", Source: "source2"})

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notification?source=source2", server.URL), nil)
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var n []Notification
	json.NewDecoder(resp.Body).Decode(&n)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(n))
	assert.Equal(t, "Test2", n[0].Message)
}
