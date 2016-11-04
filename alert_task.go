package main

import (
	"bytes"
	"fmt"
	"log"
)

// Checks for and sends alerts
func SendAlerts() {
	log.Print("Checking for alerts...")
	var n []Notification

	App.db.
		Preload("Tags").
		Preload("Tags.AlertGroups").
		Preload("Tags.AlertGroups.Recipients").
		Where(Notification{Alerted: false, Read: false}).
		Find(&n)

	if len(n) > 0 {
		filterAndSend(n)
	}
}

func filterAndSend(notifications []Notification) {
	var s []Notification
	m := make(map[string][]Notification)
	a := make(map[string][]AlertGroup)

	for _, n := range notifications {
		// if there is no tag, skip
		if n.Tags == nil {
			s = append(s, n)
			continue
		}

		for _, v := range n.Tags {

			// If there is no alert group
			// there is nothing to do here
			if v.AlertGroups == nil {
				s = append(s, n)
				continue
			}

			var nList []Notification

			if nl, ok := m[v.Name]; ok {
				nList = nl
			}

			nList = append(nList, n)
			m[v.Name] = nList
			a[v.Name] = v.AlertGroups
		}
	}

	// Send mail for ones with alert group details
	for k, _ := range m {
		if _, ok := a[k]; ok {
			sendMail(k, m[k], a[k])
		}
	}

	// Skip the ones with no alert group
	if len(s) > 0 {
		skipNotifications(s)
	}
}

func sendMail(t string, n []Notification, a []AlertGroup) {
	m := App.mg.NewMessage("notify@valeska.co.uk", fmt.Sprintf("New notifications for %s", t), "New Notifications!")
	m.SetHtml(createEmailBody(n, t))

	for _, v := range a {
		for _, r := range v.Recipients {
			m.AddRecipient(r.Email)
		}
	}

	if !App.test {
		_, _, err := App.mg.Send(m)

		if err != nil {
			// We don't want to update the notifications
			// but we also don't want to kill the server
			// with a panic.
			log.Print(err)
			return
		}
	}

	updateNotifications(n)
}

func createEmailBody(n []Notification, t string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("<html><div><p>New notifications for the tag %s</p><ul>", t))

	for _, v := range n {
		buf.WriteString(fmt.Sprintf("<li>%s (%s)- %s</li>", v.Message, v.Source, v.Action))
	}

	buf.WriteString("</ul><p>Do not reply to this email</p></div></html>")
	return buf.String()
}

func skipNotifications(n []Notification) {
	tx := App.db.Begin()

	for _, v := range n {
		v.Read = true
		if err := tx.Save(v).Error; err != nil {
			log.Print(err)
			tx.Rollback()
		}
	}

	tx.Commit()
}

func updateNotifications(n []Notification) {
	tx := App.db.Begin()

	for _, v := range n {
		v.Alerted = true
		v.Read = true
		if err := tx.Save(&v).Error; err != nil {
			log.Print(err)
			tx.Rollback()
		}
	}

	tx.Commit()
}
