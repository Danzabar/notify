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

func TestItShouldUpdateNotificationsOnSend(t *testing.T) {
    a := AlertGroup{Name: "NewTestGroup", Type: "email", Emails: "danzabian@gmail.com"}
    App.db.Create(&a)

    f := Tag{Name: "test"}
    f.AlertGroups = append(f.AlertGroups, a)
    App.db.Create(&f)

    n := Notification{Message: "Test message"}
    n.Tags = append(n.Tags, f)
    App.db.Create(&n)

    App.mg = &MockAlerter{Pass: true}
    SendAlerts()

    var o Notification
    App.db.Where("ext_id = ?", n.ExtId).Find(&o)

    assert.Equal(t, true, o.Alerted)
    assert.Equal(t, true, o.Read)
}

func TestItShouldPickUpNotificationsThatHaveAlreadyBeenAlerted(t *testing.T) {
    a := AlertGroup{Name: "NewTestGroup", Type: "email", Emails: "danzabian@gmail.com"}
    App.db.Create(&a)

    f := Tag{Name: "test"}
    f.AlertGroups = append(f.AlertGroups, a)
    App.db.Create(&f)

    n := Notification{Message: "Test message", Alerted: true, Read: false}
    n.Tags = append(n.Tags, f)
    App.db.Create(&n)

    App.mg = &MockAlerter{Pass: true}
    SendAlerts()

    var o Notification
    App.db.Where("ext_id = ?", n.ExtId).Find(&o)

    assert.Equal(t, false, o.Read)
}
