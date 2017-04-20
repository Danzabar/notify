package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	App = NewApp(":8888", "sqlite3", "/tmp/test.db", "test", "test")

	Migrate()

	App.db.Delete(&Notification{})
	App.db.Delete(&AlertGroup{})
	App.db.Delete(&Tag{})
}
