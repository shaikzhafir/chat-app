package models

type Message struct {
	User      string `json:"user"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}
