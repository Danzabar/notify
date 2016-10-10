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
	App = NewApp(":5000", "sqlite3", "/tmp/test.db")
	App.setRoutes()

	Migrate()
	server = httptest.NewServer(App.router)
}

func TestPostAlertGroupSuccess(t *testing.T) {
	payload := []byte(`{"name": "test","type":"mail"}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var a AlertGroup
	err = json.NewDecoder(resp.Body).Decode(&a)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test", a.Name)
	assert.Equal(t, "mail", a.Type)
}

func TestPostAlertGroupBadJSON(t *testing.T) {
	payload := byte(`{"test":}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPostAlertGroupValidationError(t *testing.T) {
	payload := []byte(`{"type": "test"}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPostAlertGroupUnProcessable(t *testing.T) {

}
