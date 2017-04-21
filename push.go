package main

import (
    pb "github.com/mitsuse/pushbullet-go"
    "github.com/mitsuse/pushbullet-go/requests"
    "os"
)

type PushBullet struct {
    Client *pb.Pushbullet
}

func NewPushBullet() *PushBullet {
    return &PushBullet{
        Client: pb.New(os.Getenv("PB_Token")),
    }
}

func (p *PushBullet) SendNotification(a AlertGroup, t string) bool {
    n := requests.NewNote()
    n.Title = "New Notifications"
    n.Body = t

    if _, err := p.Client.PostPushesNote(n); err != nil {
        App.log.Error(err)
        return false
    }

    return true
}
