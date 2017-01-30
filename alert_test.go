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
	App = NewApp(":5000", "sqlite3", "/tmp/test.db", "test", "test")
	App.setRoutes()

	Migrate()
	server = httptest.NewServer(App.router)
}

func TestPostAlertGroupSuccess(t *testing.T) {
	App.db.Unscoped().Delete(&AlertGroup{})
	r := &AlertGroupRequest{
		Group: AlertGroup{
			Name: "test",
			Type: "mail",
		},
	}
	payload, _ := json.Marshal(r)

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
	assert.Equal(t, "test", a.Name)
	assert.Equal(t, "mail", a.Type)
}

func TestPostAlertGroupBadJSON(t *testing.T) {
	payload := []byte(`{"test":}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPostAlertGroupValidationError(t *testing.T) {
	App.db.Unscoped().Delete(&AlertGroup{})
	r := &AlertGroupRequest{
		Group: AlertGroup{
			Type: "mail",
		},
	}

	payload, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPostAlertGroupUnProcessable(t *testing.T) {
	App.db.Delete(&AlertGroup{})
	App.db.Create(&AlertGroup{Name: "Test-Group"})
	r := &AlertGroupRequest{
		Group: AlertGroup{
			Name: "Test-Group",
			Type: "mail",
		},
	}
	payload, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/alert-group", server.URL), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 422, resp.StatusCode)
}

func TestGetAlertGroup(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/alert-group", server.URL), nil)
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestPutAlertGroupSuccess(t *testing.T) {
	a := &AlertGroup{Name: "Test Group", Type: "urgent"}
	App.db.Create(a)
	r := &AlertGroupRequest{
		Group: AlertGroup{
			Name: "updated name",
		},
	}
	payload, _ := json.Marshal(r)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/alert-group/%s", server.URL, a.ExtId), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	var g AlertGroup
	json.NewDecoder(resp.Body).Decode(&g)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "updated name", g.Name)
	assert.Equal(t, "urgent", g.Type)
}

func TestPutAlertGroupBadJSON(t *testing.T) {
	a := &AlertGroup{Name: "JSONFail", Type: "Urgent"}
	App.db.Create(a)
	payload := []byte(`{"name":}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/alert-group/%s", server.URL, a.ExtId), bytes.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 400, resp.StatusCode)
}

func TestPutAlertGroupNotFound(t *testing.T) {
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/alert-group/test-1234", server.URL), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 404, resp.StatusCode)
}

func TestPutAlertGroupUnsaveable(t *testing.T) {
	a := &AlertGroup{Name: "NameFail", Type: "Urgent"}
	App.db.Create(a)
	App.db.Create(&AlertGroup{Name: "Conflict"})

	r := &AlertGroupRequest{
		Group: AlertGroup{
			Name: "Conflict",
		},
	}
	payload, _ := json.Marshal(r)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/alert-group/%s", server.URL, a.ExtId), bytes.NewReader(payload))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 422, resp.StatusCode)
}

func TestDeleteAlertGroupSuccess(t *testing.T) {
	a := &AlertGroup{Name: "TestDelete"}
	App.db.Create(a)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/alert-group/%s", server.URL, a.ExtId), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 202, resp.StatusCode)
}

func TestDeleteAlertGroupNotFound(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/alert-group/test-1234", server.URL), nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 404, resp.StatusCode)
}
