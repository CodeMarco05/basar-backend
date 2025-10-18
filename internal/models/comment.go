package models

import "time"

type Comment struct {
	ID            string    `bson:"id" json:"id"`
	Message       string    `bson:"message" json:"message" validate:"required"`
	CommenterMail string    `bson:"commenterMail" json:"commenterMail" validate:"required"`
	CommenterName string    `bson:"commenterName" json:"commenterName" validate:"required"`
	LikedCounter  int       `bson:"likedCounter" json:"likedCounter"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
}

type CommentInsert struct {
	Message       string    `bson:"message" json:"message" validate:"required"`
	CommenterMail string    `bson:"commenterMail" json:"commenterMail" validate:"required"`
	CommenterName string    `bson:"commenterName" json:"commenterName" validate:"required"`
	LikedCounter  int       `bson:"likedCounter" json:"likedCounter"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
}
