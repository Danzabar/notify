package main

import (
	"github.com/jinzhu/gorm"
)

// Tag Model
type Tag struct {
	gorm.Model

	Name string
}

// Notification Model
type Notification struct {
	gorm.Model

	Message string `gorm:"type:text"`
	Tags    []Tag
}
