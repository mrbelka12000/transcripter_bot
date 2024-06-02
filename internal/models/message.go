package models

type Message struct {
	Text      string `bson:"text,omitempty"`
	MessageID int64  `bson:"message_id,omitempty"`
	ChatID    int64  `bson:"chat_id,omitempty"`
}
