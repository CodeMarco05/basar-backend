package models

import (
	"time"
)

type Post struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatorId   string    `bson:"creatorId" json:"creatorId"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Tags        []string  `bson:"tags" json:"tags"`
	Text        string    `bson:"text" json:"text"`
	Mail        string    `bson:"mail" json:"mail"`
	PayPalMail  string    `bson:"payPalMail" json:"payPalMail"`
	Images      []string  `bson:"images" json:"images"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type InsertPost struct {
	CreatorId   string   `bson:"creatorId" json:"creatorId" validate:"required"`
	Title       string   `bson:"title" json:"title" validate:"required"`
	Description string   `bson:"description" json:"description" validate:"required"`
	Tags        []string `bson:"tags" json:"tags" validate:"required"`
	Text        string   `bson:"text" json:"text" validate:"required"`
	Mail        string   `bson:"mail" json:"mail" validate:"required"`
	PayPalMail  string   `bson:"payPalMail" json:"payPalMail" validate:"required,email"`
	Images      []string `bson:"images" json:"images" validate:"required"`
}
