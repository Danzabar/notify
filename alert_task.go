package main

// Checks for and sends alerts
func SendAlerts() {
	var n []Notification

	App.db.Preload("Tags").
		Preload("Tags.AlertGroups").
		Where(Notification{Alerted: false, Read: false}).
		Find(&n)
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
}

func SendPushNotification(a AlertGroup, n []Notification) bool {
	return true
}

func SendEmailNotification(a AlertGroup, n []Notification) bool {
	return true
}
