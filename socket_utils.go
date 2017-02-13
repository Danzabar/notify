package main

import (
    "reflect"
)

func SearchNotifications(r NotificationRefresh, p *Pagination) (c int, n []Notification) {
    s := GetSearchOptionsFromRequest(r)

    App.db.Model(&Notification{}).Where(s).Count(&c)
    App.db.Model(&Notification{}).
        Where(s).
        Preload("Tags").
        Order("updated_at DESC").
        Limit(p.Limit).
        Offset(p.Offset).
        Find(&n)

    return
}

func GetSearchOptionsFromRequest(r NotificationRefresh) Notification {
    n := Notification{}

    if r.Read {
        n.Read = true
    }

    if len(r.Tags) > 0 {
        for _, v := range r.Tags {
            n.Tags = append(n.Tags, Tag{Name: v})
        }
    }

    if len(r.Search) > 0 {
        ref := reflect.ValueOf(n)
        s := ref.Elem()

        for k, v := range r.Search {
            f := s.FieldByName(k)

            if f.IsValid() && f.CanSet() {
                f.SetString(v)
            }
        }
    }

    return n
}
