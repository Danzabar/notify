package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	App = NewApp(":5000", "sqlite3", "/tmp/test.db")

	Migrate()
}

func TestSocketLoadEvent(t *testing.T) {
	var p SocketLoadPayload
	App.db.Create(&Notification{Message: "Test"})

	resp := App.OnSocketLoad()
	err := json.NewDecoder(bytes.NewReader(resp)).Decode(&p)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(p.Notifications) > 0)
}
