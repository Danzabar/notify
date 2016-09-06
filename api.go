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
