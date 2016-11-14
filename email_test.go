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

func TestGetEmailHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/email", server.URL), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostEmailSuccess(t *testing.T) {
	r := &EmailTemplateRequest{}
	r.Template = EmailTemplate{Name: "test", Content: "test"}
	r.Tags = append(r.Tags, Tag{Name: "test-tag", Source: "Test"})

	reqPayload, _ := json.Marshal(r)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/email", server.URL), bytes.NewReader(reqPayload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var e EmailTemplateRequest
	json.NewDecoder(resp.Body).Decode(&e)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test", e.Template.Name)
	assert.Equal(t, 1, len(e.Tags))
}
