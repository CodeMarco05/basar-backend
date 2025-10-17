package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Creator     string             `bson:"creator" json:"creator"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Tags        []string           `bson:"tags" json:"tags"`
	Text        string             `bson:"text" json:"text"`
	PayPalMail  string             `bson:"payPalMail" json:"payPalMail"`
	Images      []string           `bson:"images" json:"images"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
