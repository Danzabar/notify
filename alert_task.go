package main

import (
	"fmt"
)

// Checks for and sends alerts
func SendAlerts() {
	var n []Notification

	App.log.Debug("Running Alert process...")

	App.db.Where("alerted != ?", true).
		Preload("Tags").
		Preload("Tags.AlertGroups").
		Find(&n)

	FilterCorrespondence(n)

	App.log.Debug("Finished")
}

func FilterCorrespondence(n []Notification) {
	var s []Notification

	App.log.Debugf("Found %d notifications", len(n))

	for _, v := range n {
		for _, t := range v.Tags {

			if len(t.AlertGroups) == 0 {
				s = append(s, v)
				continue
			}

			for _, a := range t.AlertGroups {
				switch a.Type {
				case "email":
					if App.mg.SendNotification(a, createEmailBody(v, t.Name)) {
						s = append(s, v)
					}
					break
				case "push":
					msg := fmt.Sprintf("%s\n%s\nFrom: %s", v.Message, v.Action, v.Source)
					if App.pb.SendNotification(a, msg) {
						s = append(s, v)
					}
					break
				}
			}
		}
	}

	UpdateNotifications(s)
}

func UpdateNotifications(n []Notification) {
	App.log.Debug("Finishing alerting process...")
	for _, v := range n {
		v.Alerted = true
		v.Read = true

		App.db.Save(&v)
	}
}

func createEmailBody(n Notification, t string) string {
	return fmt.Sprintf("<html><div><p>New notification for %s</p><p>%s</p><p>%s</p><p>From: %s</p><p>Do not reply to this email</p></div></html>", t, n.Message, n.Action, n.Source)
}
