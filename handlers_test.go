package main

import (
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

func TestPingHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/ping", server.URL), nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode)
}

func TestPaginationOptions(t *testing.T) {
	req1, _ := http.NewRequest("GET", "/test?pageSize=4&page=3", nil)
	req2, _ := http.NewRequest("GET", "/test?pageSize=1&page=1", nil)
	req3, _ := http.NewRequest("GET", "/test", nil)

	p1 := GetPaginationFromRequest(req1)
	p2 := GetPaginationFromRequest(req2)
	p3 := GetPaginationFromRequest(req3)

	assert.Equal(t, 50, p3.Limit)
	assert.Equal(t, 0, p3.Offset)

	assert.Equal(t, 1, p2.Limit)
	assert.Equal(t, 0, p2.Offset)

	assert.Equal(t, 4, p1.Limit)
	assert.Equal(t, 8, p1.Offset)
}
