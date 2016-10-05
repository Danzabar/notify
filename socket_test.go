package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func init() {
	App = NewApp(":5000", "sqlite3", "/tmp/test.db")

	Migrate()
}

func TestSocketReadEvent(t *testing.T) {
	n := &Notification{Message: "A test read event"}
	App.db.Create(n)

	payload := fmt.Sprintf(`{"ids":["%s"]}`, n.ExtId)

	App.OnNotificationRead(payload)

	var out Notification
	App.db.Where(&Notification{ExtId: n.ExtId}).First(&out)

	assert.Equal(t, true, out.Read)
}

func TestSocketReadEventBadJson(t *testing.T) {
	var r Response
	payload := `{"ids":}`

	resp := App.OnNotificationRead(payload)

	err := json.NewDecoder(strings.NewReader(resp)).Decode(&r)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Invalid json", r.Error)
}

func TestSocketRefreshEvent(t *testing.T) {
	var r SocketLoadPayload
	App.db.Delete(&Notification{})
	App.db.Create(&Notification{Message: "Test Refresh"})
	payload := `{"page":0, "pageSize": 10}`

	resp := App.OnNotificationRefresh(payload)

	err := json.NewDecoder(strings.NewReader(resp)).Decode(&r)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(r.Notifications))
	assert.Equal(t, false, r.HasNext)
	assert.Equal(t, false, r.HasPrev)
}

func TestSocketRefreshEventWithNext(t *testing.T) {
	var r SocketLoadPayload
	App.db.Delete(&Notification{})
	App.db.Create(&Notification{Message: "Test1"})
	App.db.Create(&Notification{Message: "Test2"})
	payload := `{"page": 1, "pageSize": 1}`

	resp := App.OnNotificationRefresh(payload)
	err := json.NewDecoder(strings.NewReader(resp)).Decode(&r)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(r.Notifications))
	assert.Equal(t, true, r.HasNext)
	assert.Equal(t, false, r.HasPrev)
}

func TestSocketRefreshEventWithPrev(t *testing.T) {
	var r SocketLoadPayload
	App.db.Delete(&Notification{})
	App.db.Create(&Notification{Message: "Test1"})
	App.db.Create(&Notification{Message: "Test2"})
	App.db.Create(&Notification{Message: "Test3"})
	App.db.Create(&Notification{Message: "Test4"})
	payload := `{"page":2, "pageSize": 2}`

	resp := App.OnNotificationRefresh(payload)
	err := json.NewDecoder(strings.NewReader(resp)).Decode(&r)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(r.Notifications))
	assert.Equal(t, false, r.HasNext)
	assert.Equal(t, true, r.HasPrev)
}
