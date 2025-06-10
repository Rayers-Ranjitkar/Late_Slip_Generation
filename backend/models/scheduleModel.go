package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: need to update the model
type Schedule struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	ModuleCode     string             `bson:"module_code" json:"module_code"`
	ModuleName     string             `bson:"module_name" json:"module_name"`
	StartTime      string             `bson:"start_time" json:"start_time"` // string for now
	EndTime        string             `bson:"end_time" json:"end_time"`     // string for now
	Day            string             `bson:"day" json:"day"`
	RoomName       string             `bson:"room_name" json:"room_name"`
	InstructorName string             `bson:"instructor_name" json:"instructor_name"`
	Semester       string             `bson:"semester" json:"semester"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
