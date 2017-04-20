package main

import (
	"encoding/json"
	"net/http"
)

// Rest Response contract
type RestResponse interface {
	Serialize() []byte
}

type RestRequest interface {
	Deserialize(r *http.Request) error
}

// Request Body for Notification Endpoint
type NotificationRequest struct {
	Notifications []Notification `json:"notifications"`
	Tags          []Tag          `json:"tags"`
}

func (n *NotificationRequest) Deserialize(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(n)
}

// Request Body for Alert groups
type AlertGroupRequest struct {
	Group AlertGroup `json:"group"`
	Tags  []Tag      `json:"tags"`
}

func (a *AlertGroupRequest) Deserialize(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(a)
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
	Notifications []Notification `json:"notifications"`
	Tags          []Tag          `json:"tags"`
	HasNext       bool           `json:"hasNext"`
	HasPrev       bool           `json:"hasPrev"`
}

func (s *SocketLoadPayload) Serialize() []byte {
	jsonResp, _ := json.Marshal(s)
	return jsonResp
}

// Request sent from client socket to tell us which notifications have been read
type NotificationRead struct {
	Ids []string `json:"ids"`
}

// Request sent from client socket to refresh notifications by page
type NotificationRefresh struct {
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Tags     []string          `json:"tags"`
	Read     bool              `json:"read"`
	Search   map[string]string `json:"search"`
}
