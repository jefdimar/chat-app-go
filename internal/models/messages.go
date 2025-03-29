package models

// Message represents a chat message
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}
