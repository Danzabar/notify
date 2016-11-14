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
	r.Template = EmailTemplate{Name: "Test"}
	r.Tags = append(r.Tags, Tag{Name: "test-tag", Source: "Test"})

	reqPayload, err := json.Marshal(r)

	if err != nil {
		t.Fatal(err)
	}

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

func TestPostEmailBadJson(t *testing.T) {
	reqPayload := []byte(`{"test":}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/email", server.URL), bytes.NewReader(reqPayload))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPostEmailValidationError(t *testing.T) {
	r := &EmailTemplateRequest{}
	r.Template = EmailTemplate{Content: "Some Content"}

	reqPayload, _ := json.Marshal(r)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/email", server.URL), bytes.NewReader(reqPayload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestFindEmailTemplateSuccess(t *testing.T) {
	e := &EmailTemplate{Name: "Test", Content: "Test Content"}
	App.db.Create(e)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/email/%s", server.URL, e.ExtId), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestFuncEmailTemplateNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/email/test", server.URL), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 404, resp.StatusCode)
}
