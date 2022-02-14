package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	BaseId        `bson:",inline"`
	BaseCreatable `bson:",inline"`

	Activity string   `bson:"activity" json:"activity"`
	Metadata bson.M   `bson:"metadata" json:"metadata"`
	Payload  []bson.M `bson:"payload" json:"payload"`
}

func MakeEvent() Event {
	o := Event{}
	o.Id = primitive.NewObjectID()
	o.CreatedAt = time.Now()

	return o
}
