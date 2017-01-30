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

func init() {
	App = NewApp("5000", "sqlite3", "/tmp/test.db", "test", "test")
	App.setRoutes()

	Migrate()

	server = httptest.NewServer(App.router)
}

// Test that Get Tag endpoint returns 200
func TestGetTagReturns200(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tag", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/tag")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// Test post endpoint success
func TestPostTagSuccess(t *testing.T) {
	reqPayload := []byte(`{"name":"test1"}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tag", server.URL), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/tag")
	}

	resp, _ := http.DefaultClient.Do(req)

	var tag Tag

	json.NewDecoder(resp.Body).Decode(&tag)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test", tag.Name)
}

func TestPostTagValidation(t *testing.T) {
	reqPayload := []byte(`{"name":""}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tag", server.URL), bytes.NewReader(reqPayload))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

// Test that bad json returns 400
func TestPostTagBadJson(t *testing.T) {
	reqPayload := []byte(`{"test":}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tag", server.URL), bytes.NewReader(reqPayload))

	if err != nil {
		t.Fatal("Request failed [POST] /api/v1/tag")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestFindTagSuccess(t *testing.T) {
	tag := &Tag{Name: "findme"}
	App.db.Create(tag)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tag/%s", server.URL, tag.ExtId), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/tag/{id}")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestFindTagNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tag/not-found", server.URL), nil)

	if err != nil {
		t.Fatal("Request failed [GET] /api/v1/tag/{id}")
	}

	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}
