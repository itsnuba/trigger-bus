package requests

import (
	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
)

type EventAddConverter interface {
	FromEvent(i models.Event) error
	ToEvent() (models.Event, error)
}

type EventAddForm struct {
	Activity string   `json:"activity" binding:"required,activityFormat"`
	Metadata *bson.M  `json:"metadata" binding:"required"`
	Payload  []bson.M `json:"payload" binding:"required"`
}

type EventAddFormMetadata struct {
	AccountId      int    `json:"accountId" binding:"gte=0"`
	InvokerId      int    `json:"invokerId" binding:"gte=0"`
	MicroserviceId string `json:"microserviceId" binding:"required"`
}

func (f *EventAddForm) FromEvent(i models.Event) error {
	f.Activity = i.Activity

	// metadata
	f.Metadata = &i.Metadata

	// payload
	f.Payload = i.Payload

	return nil
}

func (f EventAddForm) ToEvent() (models.Event, error) {
	o := models.MakeEvent()
	o.Activity = f.Activity
	o.Metadata = bson.M{}
	o.Payload = []bson.M{}

	// metadata
	if f.Metadata != nil {
		o.Metadata = *f.Metadata
	}

	// payload
	if f.Payload != nil {
		o.Payload = f.Payload
	}

	return o, nil
}
