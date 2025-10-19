package models

import (
	"time"
)

type Post struct {
	ID             string              `bson:"_id,omitempty" json:"id,omitempty"`
	CreatorId      string              `bson:"creatorId" json:"creatorId"`
	CreatorMail    string              `bson:"creatorMail" json:"creatorMail" validate:"required"`
	Title          string              `bson:"title" json:"title" validate:"required"`
	Description    string              `bson:"description" json:"description" validate:"required"`
	Tags           []string            `bson:"tags" json:"tags" validate:"required"`
	Text           string              `bson:"text" json:"text" validate:"required"`
	PayPalMail     string              `bson:"payPalMail" json:"payPalMail"`
	Images         []string            `bson:"images" json:"images" validate:"required"`
	Comments       []Comment           `bson:"comments" json:"comments" validate:"required"`
	IsTerminated   *bool               `bson:"isTerminated" json:"isTerminated" validate:"required"`
	AcceptanceList []AcceptanceRequest `bson:"acceptanceList" json:"acceptanceList" validate:"required"`
	AcceptedUser   AcceptedUser        `bson:"acceptedUser,omitempty" json:"acceptedUser,omitempty" validate:"required"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
}

type InsertPost struct {
	CreatorId   string    `bson:"creatorId" json:"creatorId" validate:"required"`
	CreatorMail string    `bson:"creatorMail" json:"creatorMail" validate:"required"`
	Title       string    `bson:"title" json:"title" validate:"required"`
	Description string    `bson:"description" json:"description" validate:"required"`
	Tags        []string  `bson:"tags" json:"tags" validate:"required"`
	Text        string    `bson:"text" json:"text" validate:"required"`
	PayPalMail  string    `bson:"payPalMail" json:"payPalMail" validate:"omitempty"`
	Comments    []Comment `bson:"comments" json:"comments"`
	Images      []string  `bson:"images" json:"images" validate:"required"`
}

type AcceptanceRequest struct {
	UserId      string `json:"userId" validate:"required"`
	UserName    string `json:"userName" validate:"required"`
	RequestedAt string `bson:"requestedAt" json:"requestedAt"`
}

type AcceptedUser struct {
	UserId     string `json:"userId" validate:"omitempty"`
	UserName   string `json:"userName" validate:"omitempty"`
	AcceptedAt string `json:"acceptedAt"`
}
