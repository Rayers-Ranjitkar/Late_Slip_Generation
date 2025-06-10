package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Fullname  string            `bson:"fullname" json:"fullname" binding:"required"`
    Password  string            `bson:"password" json:"password" binding:"required"`  
    Email     string            `bson:"email" json:"email" binding:"required,email"`
    Role      string            `bson:"role" json:"role" binding:"oneof=student admin"`
    CreatedAt time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}
