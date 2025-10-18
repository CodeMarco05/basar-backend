package models

import "time"

type Comment struct {
	Message       string    `bson:"message" json:"message" validate:"required"`
	CommenterMail string    `bson:"commenterMail" json:"commenterMail" validate:"required"`
	CommenterName string    `bson:"commenterName" json:"commenterName" validate:"required"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
}
