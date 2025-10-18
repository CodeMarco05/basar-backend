package models

import (
	"time"
)

type Post struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatorId   string    `bson:"creatorId" json:"creator"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Tags        []string  `bson:"tags" json:"tags"`
	Text        string    `bson:"text" json:"text"`
	PayPalMail  string    `bson:"payPalMail" json:"payPalMail"`
	Images      []string  `bson:"images" json:"images"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type InsertPost struct {
	CreatorId   string   `json:"creatorId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Text        string   `json:"text"`
	PayPalMail  string   `json:"payPalMail"`
	Images      []string `json:"images"`
}
