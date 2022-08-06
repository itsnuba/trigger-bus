package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TriggerScheduler struct {
	BaseId        `bson:",inline"`
	BaseCreatable `bson:",inline"`
	BaseModifable `bson:",inline"`

	Metadata bson.M `bson:"metadata" json:"metadata"`
	CronExpr string `bson:"cronExpr" json:"cronExpr"`

	Active      bool   `bson:"active" json:"active"`
	EndpointUrl string `bson:"endpointUrl" json:"endpointUrl"`
}

func MakeTriggerScheduler() TriggerScheduler {
	o := TriggerScheduler{
		Metadata: bson.M{},
	}
	o.Id = primitive.NewObjectID()
	o.CreatedAt = time.Now()

	return o
}

func (v *TriggerScheduler) ToTriggerLog() TriggerLog {
	o := MakeTriggerLog()
	o.TriggerType = TLETScheduler
	o.TriggerListenerId = &v.Id
	o.CallbackUrl = v.EndpointUrl

	return o
}
