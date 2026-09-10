package model

type APIResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Type    string      `json:"type,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
