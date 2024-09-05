package models

type Message struct {
	Text      string `bson:"text,omitempty"`
	MessageID int    `bson:"message_id,omitempty"`
	ChatID    string `bson:"chat_id,omitempty"`
}
