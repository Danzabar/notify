package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	App = NewApp(":8888", "sqlite3", "/tmp/test.db")

	Migrate()

	App.db.Delete(&Notification{})
	App.db.Delete(&AlertGroup{})
	App.db.Delete(&Tag{})
}

func TestPendingAlertsAreUpdated(t *testing.T) {
	n := &Notification{Message: "Test Alerting", Source: "Test"}
	m := &Tag{
		Name: "TestAlert",
		AlertGroups: []AlertGroup{
			AlertGroup{
				Name: "TestGroup",
				Recipients: []Recipient{
					Recipient{
						Email: "danzabian@gmail.com",
					},
				},
			},
		},
	}

	App.db.Create(m)
	App.db.Create(n)
	App.db.Model(n).Association("Tags").Append(m)

	App.test = true
	SendAlerts()

	var u Notification
	App.db.Where("ext_id = ?", n.ExtId).First(&u)

	assert.Equal(t, true, u.Alerted)
}

func TestSkipRecordsWithNoAlertGroup(t *testing.T) {
	n := &Notification{Message: "Test Skip Alert", Source: "Test"}
	m := &Tag{Name: "TestSkipAlert"}

	App.db.Create(n)
	App.db.Create(m)
	App.db.Model(n).Association("Tags").Append(m)

	App.test = true
	SendAlerts()

	var u Notification
	App.db.Where("ext_id = ?", n.ExtId).First(&u)

	assert.Equal(t, true, u.Read)
	assert.Equal(t, false, u.Alerted)
}
