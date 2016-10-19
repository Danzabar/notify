package main

import (
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
		Where("alerted = ?", false).
		Where("read = ?", false).
		Find(&n)

	if len(n) > 0 {
		log.Printf("Found %d notifications to alert on", len(n))
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
		log.Printf("Skipping %d notifications", len(s))
		updateNotifications(s)
	}
}

func sendMail(t string, n []Notification, a []AlertGroup) error {
	m := App.mg.NewMessage("notify@valeska.co.uk", fmt.Sprintf("New notifications for %s", t), "You have notifications")

	for _, v := range a {
		for _, r := range v.Recipients {
			log.Printf("adding %s recipient", r.Email)
			m.AddRecipient(r.Email)
		}
	}

	_, _, err := App.mg.Send(m)

	if err != nil {
		// We don't want to update the notifications
		// but we also don't want to kill the server
		// with a panic.
		log.Fatal(err)
		return
	}

	updateNotifications(n)
}

func createEmailBody(n []Notification) {

}

func updateNotifications(n []Notification) {
	tx := App.db.Begin()

	for _, v := range n {
		v.Alerted = true
		if err := tx.Save(&v).Error; err != nil {
			tx.Rollback()
			updateNotifications(n)
		}
	}

	tx.Commit()
	log.Printf("Updated %d notifications", len(n))
}
