package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TriggerLog struct {
	BaseId        `bson:",inline"`
	BaseCreatable `bson:",inline"`
	BaseModifable `bson:",inline"`

	// listener
	TriggerListenerId *primitive.ObjectID `bson:"_triggerListenerId" json:"triggerListenerId"`
	CallbackUrl       string              `bson:"callbackUrl" json:"callbackUrl"`
	HandlingParameter bson.M              `bson:"handlingParameter" json:"handlingParameter"`
	Handled           bool                `bson:"handled" json:"handled"`
	HandledMessage    *string             `bson:"handledMessage" json:"handledMessage"`

	// event
	EventId primitive.ObjectID `bson:"_eventId" json:"eventId"`
}

func MakeTriggerLog() TriggerLog {
	o := TriggerLog{}
	o.Id = primitive.NewObjectID()
	o.CreatedAt = time.Now()

	return o
}
