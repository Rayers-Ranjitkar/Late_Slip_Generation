package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LateSlip struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID primitive.ObjectID `bson:"student_id" json:"student_id"`
	Reason    string             `bson:"reason" json:"reason" binding:"required"`
	Status    string             `bson:"status" json:"status" binding:"required,oneof=pending approved rejected"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	RequestID string             `bson:"request_id,omitempty" json:"request_id"`
}
