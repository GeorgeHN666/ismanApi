package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriber struct {
	ID        primitive.ObjectID `json:"_id"  bson:"_id"`
	FullName  string             `json:"full_name"  bson:"full_name"`
	Email     string             `json:"email"  bson:"email"`
	Code      string             `json:"code"  bson:"code"`
	Status    int                `json:"status"  bson:"status"`
	CreatedAt time.Time          `json:"created_at"   bson:"created_at"`
}
