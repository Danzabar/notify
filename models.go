package main

import (
	"encoding/json"
	"github.com/influxdata/influxdb/uuid"
	"github.com/jinzhu/gorm"
)

// Tag Model
type Tag struct {
	gorm.Model

	// The Tag Name
	Name string `json:"name"`
	// The source system
	Source string `json:"source"`
	// External ID
	ExtId string `json:"extid"`
}

// Tag Before Create
func (t *Tag) BeforeCreate() {
	// Create UUID On create
	u := uuid.TimeUUID()
	t.ExtId = u.String()
}

func (t *Tag) AfterCreate() {
	jsonStr, _ := json.Marshal(t)

	App.server.BroadcastTo("notify", "new:tag", string(jsonStr))
}

// Notification Model
type Notification struct {
	gorm.Model

	// The message of the notification
	Message string `gorm:"type:text" json:"message"`
	// Any action that should be taken on this notification
	Action string `json:"action"`
	// The External ID
	ExtId string `json:"extid"`
	// Source system
	Source string `json:"source"`
	// Flag for read
	Read bool `json:"read"`
	// List of related Tags
	Tags []Tag `json:"tags,omitempty"`
}

// Notification Before Create
func (n *Notification) BeforeCreate() {
	// Create UUID on create
	u := uuid.TimeUUID()
	n.ExtId = u.String()
}

func (n *Notification) AfterCreate() {
	jsonStr, _ := json.Marshal(n)

	// Broadcast creation to socket listeners
	App.server.BroadcastTo("notify", "new:notify", string(jsonStr))
}
