package models

type Item struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}
