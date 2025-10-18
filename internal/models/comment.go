package models

import "time"

type Comment struct {
	Message       string    `bson:"message" json:"message" validate:"required"`
	CommenterMail string    `bson:"mail" json:"CommenterMail" validate:"required"`
	CommenterName string    `bson:"name" json:"CommenterName" validate:"required"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
}
