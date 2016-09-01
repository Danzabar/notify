package main

import (
	"github.com/influxdata/influxdb/uuid"
	"github.com/jinzhu/gorm"
)

// Tag Model
type Tag struct {
	gorm.Model `json:"-"`

	Name  string `json:"name"`
	ExtId string `json:"id"`
}

// Tag Before Create
func (t *Tag) BeforeCreate() {
	u := uuid.TimeUUID()
	t.ExtId = u.String()
}

// Notification Model
type Notification struct {
	gorm.Model `json:"-"`

	Message string `gorm:"type:text" json:"message"`
	ExtId   string `json:"id"`
	Tags    []Tag  `json:"tags,omitempty"`
}

// Notification Before Create
func (n *Notification) BeforeCreate() {
	u := uuid.TimeUUID()
	n.ExtId = u.String()
}
