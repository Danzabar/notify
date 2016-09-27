package main

import (
	"encoding/json"
)

// Rest Response contract
type RestResponse interface {
	Serialize() []byte
}

// Standard response struct
type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Serializes the data in the Response struct
func (r *Response) Serialize() []byte {
	jsonResp, _ := json.Marshal(r)
	return jsonResp
}

type ValidationResponse struct {
	Errors map[string]string `json:"errors"`
}

func (v *ValidationResponse) Serialize() []byte {
	jsonResp, _ := json.Marshal(v)
	return jsonResp
}

// Payload for the socket load event
type SocketLoadPayload struct {
	Notifications []Notification `json:"notifications,omitempty"`
	Tags          []Tag          `json:"tags,omitempty"`
}

func (s *SocketLoadPayload) Serialize() []byte {
	jsonResp, _ := json.Marshal(s)
	return jsonResp
}

// Request sent from client socket to tell us which notifications have been read
type NotificationRead struct {
	Ids []string `json:"ids"`
}
