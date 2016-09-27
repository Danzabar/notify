package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSocketReadEvent(t *testing.T) {
	n := &Notification{Message: "A test read event"}
	App.db.Create(n)

	payload := []byte(fmt.Sprintf(`{"ids":["%s"]}`, n.ExtId))

	App.OnNotificationRead(string(payload))

	var out Notification
	App.db.Where(&Notification{ExtId: n.ExtId}).First(&out)

	assert.Equal(t, true, out.Read)
}
