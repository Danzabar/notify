package main

import (
	"encoding/json"
	"github.com/influxdata/influxdb/uuid"
	"github.com/leebenson/conform"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

// Tag Model
type Tag struct {
	Model

	// The Tag Name
	Name string `gorm:"not null" json:"name" conform:"name" validate:"required"`
	// The source system
	Source string `json:"source" conform:"slug"`
	// External ID
	ExtId string `gorm:"unique" json:"extid"`
	// Alert group relationship
	AlertGroups []AlertGroup `gorm:"many2many:group_tags"`
}

// Tag Before Create
func (t *Tag) BeforeCreate() {
	// Create UUID On create
	u := uuid.TimeUUID()
	t.ExtId = u.String()

	conform.Strings(t)
}

func (t *Tag) AfterCreate() {
	jsonStr, _ := json.Marshal(t)
	App.server.BroadcastTo("notify", "new:tag", string(jsonStr))
}

// Notification Model
type Notification struct {
	Model

	// The message of the notification
	Message string `gorm:"type:text;not null" json:"message" conform:"ucfirst,trim" validate:"required"`
	// Any action that should be taken on this notification
	Action string `json:"action"`
	// The External ID
	ExtId string `gorm:"unique" json:"extid"`
	// Source system
	Source string `json:"source" conform:"slug" validate:"required"`
	// Flag for read
	Read bool `json:"read"`
	// Boolean flag for when the schedule picks it up
	Alerted bool `json:"alerted"`
	// List of related Tags
	Tags []Tag `gorm:"many2many:notification_tags" json:"tags,omitempty"`
}

// Notification Before Create
func (n *Notification) BeforeCreate() {
	// Create UUID on create
	u := uuid.TimeUUID()
	n.ExtId = u.String()

	conform.Strings(n)
}

func (n *Notification) AfterCreate() {
	jsonStr, _ := json.Marshal(n)
	App.server.BroadcastTo("notify", "new:notify", string(jsonStr))
}

type AlertGroup struct {
	Model

	Name       string      `gorm:"unique" json:"name" validate:"required"`
	ExtId      string      `gorm:"unique" json:"extId"`
	Type       string      `json:"type"`
	Tags       []Tag       `gorm:"many2many:group_tags" json:"tags"`
	Recipients []Recipient `gorm:"many2many:group_recipients" json:"recipients"`
}

func (a *AlertGroup) BeforeCreate() {
	a.ExtId = uuid.TimeUUID().String()
	conform.Strings(a)
}

type Recipient struct {
	Model

	Name  string `json:"name" conform:"name"`
	Email string `json:"email"`
	ExtId string `gorm:"unique" json:"extId"`
}

func (r *Recipient) BeforeCreate() {
	r.ExtId = uuid.TimeUUID().String()
	conform.Strings(r)
}

type EmailTemplate struct {
}
