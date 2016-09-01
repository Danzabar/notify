package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
)

type RestResponse interface {
	Serialize() []byte
}

type Response struct {
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
	Data    []*gorm.Model `json:"data,omitempty"`
}

func (r *Response) Serialize() []byte {
	jsonResp, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
		return []byte("")
	}

	return jsonResp
}
