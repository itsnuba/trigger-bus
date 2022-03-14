package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TriggerListener struct {
	BaseId        `bson:",inline"`
	BaseCreatable `bson:",inline"`
	BaseModifable `bson:",inline"`

	Activity          string            `bson:"activity" json:"activity"`
	MetadataFilter    map[string]bson.A `bson:"metadataFilter" json:"metadataFilter"`
	HandlingParameter bson.M            `bson:"handlingParameter" json:"handlingParameter"`

	Active      bool   `bson:"active" json:"active"`
	CallbackUrl string `bson:"callbackUrl" json:"callbackUrl"`
}

func MakeTriggerListener() TriggerListener {
	o := TriggerListener{}
	o.Id = primitive.NewObjectID()
	o.CreatedAt = time.Now()

	return o
}

func (v *TriggerListener) ToTriggerLog() TriggerLog {
	o := MakeTriggerLog()
	o.TriggerListenerId = &v.Id
	o.CallbackUrl = v.CallbackUrl
	o.HandlingParameter = v.HandlingParameter

	return o
}
