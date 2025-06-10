package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID     string             `bson:"student_id" json:"student_id"`
	Name          string             `bson:"name" json:"name"`
	Email         string             `bson:"email" json:"email"`
	Semester      string             `bson:"semester" json:"semester"`
	Level         string             `bson:"level" json:"level"`
	LateSlipCount int                `bson:"late_slip_count" json:"late_slip_count"`
	//TODO: need to replace gender with Semester
	// -- this is just a placeholder for now
	//--- need to update the model later
}
