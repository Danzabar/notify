package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	App = NewApp("5000", "sqlite3", "/tmp/test.db")
}

func TestItReturnsLatestNotifications(t *testing.T) {
	App.db.Create(&Notification{})
	App.db.Create(&Tag{})

	resp := App.OnSocketLoad()

	var p SocketLoadPayload

	err := json.NewDecoder(bytes.NewReader(resp)).Decode(&p)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, p.Notifications)
	assert.NotEmpty(t, p.Tags)
}
