package main

import (
	"bytes"
	"fmt"
)

// Checks for and sends alerts
func SendAlerts() {
	var n []Notification

	App.db.Preload("Tags").
		Preload("Tags.AlertGroups").
		Where(Notification{Alerted: false, Read: false}).
		Find(&n)

	FilterCorrespondence(n)
}

func FilterCorrespondence(n []Notification) {
	var tl []Tag
	var s []Notification
	nm := make(map[string][]Notification)

	for _, v := range n {
		for _, t := range v.Tags {
			nm[t.ExtId] = append(nm[t.ExtId], v)
			tl = append(tl, t)
		}
	}

	for _, t := range tl {

		if len(t.AlertGroups) == 0 {
			s = append(s, nm[t.ExtId]...)
			continue
		}

		for _, a := range t.AlertGroups {
			switch a.Type {
			case "email":
				if r := SendEmailNotification(a, nm[t.ExtId]); r {
					s = append(s, nm[t.ExtId]...)
				}
				break
			case "push":
				if r := SendPushNotification(a, nm[t.ExtId]); r {
					s = append(s, nm[t.ExtId]...)
				}
				break
			}
		}
	}

	UpdateNotifications(s)
}

func UpdateNotifications(n []Notification) {
	for _, v := range n {
		v.Alerted = true
		v.Read = true

		App.db.Save(&v)
	}
}

func SendPushNotification(a AlertGroup, n []Notification) bool {
	return true
}

func SendEmailNotification(a AlertGroup, n []Notification) bool {
	return true
}

func createEmailBody(n []Notification, t string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("<html><div><p>New notifications</p><ul>", t))

	for _, v := range n {
		buf.WriteString(fmt.Sprintf("<li>%s (%s)- %s</li>", v.Message, v.Source, v.Action))
	}

	buf.WriteString("</ul><p>Do not reply to this email</p></div></html>")
	return buf.String()
}
