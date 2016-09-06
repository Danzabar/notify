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
	App = NewApp("5000", "test")
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
	assert.Equal(t, "test1", tag.Name)
}
