package main

import (
    "fmt"
    "gopkg.in/mailgun/mailgun-go.v1"
    "os"
    "strings"
)

type Mailgun struct {
    Client mailgun.Mailgun
}

func NewMailClient() *Mailgun {
    return &Mailgun{
        Client: mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_APIKEY"), os.Getenv("MG_PUBKEY")),
    }
}

func (m *Mailgun) SendNotification(a AlertGroup, t string) bool {
    ms := m.Client.NewMessage(os.Getenv("MG_FROM"), fmt.Sprintf("New notifications for %s", t), "New Notifications!")
    ms.SetHtml(t)

    for _, v := range strings.Split(a.Emails, ",") {
        ms.AddRecipient(v)
    }

    _, _, e := m.Client.Send(ms)

    if e != nil {
        App.log.Error(e)
        return false
    }

    return true
}
