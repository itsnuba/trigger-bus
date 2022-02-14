package requests

import (
	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerListenerAddConverter interface {
	FromTriggerListener(i models.TriggerListener) error
	ToTriggerListener() (models.TriggerListener, error)
}

type TriggerListenerAddForm struct {
	Activity          string                        `json:"activity" binding:"required,activityFormat"`
	CallbackUrl       string                        `json:"callbackUrl" binding:"required,url"`
	MetadataFilter    TriggerListenerMetadataFilter `json:"metadataFilter" binding:"required"`
	HandlingParameter bson.M                        `bson:"handlingParameter" binding:"required"`
}

func (f *TriggerListenerAddForm) FromTriggerListener(i models.TriggerListener) error {
	f.Activity = i.Activity
	f.CallbackUrl = i.CallbackUrl

	// metadata
	f.MetadataFilter = TriggerListenerMetadataFilter{}
	f.MetadataFilter.FromBsonM(i.MetadataFilter)

	// handling parameter
	f.HandlingParameter = i.HandlingParameter

	return nil
}

func (f TriggerListenerAddForm) ToTriggerListener() (models.TriggerListener, error) {
	o := models.MakeTriggerListener()
	o.Activity = f.Activity
	o.CallbackUrl = f.CallbackUrl
	o.HandlingParameter = f.HandlingParameter

	// metadata
	o.MetadataFilter = bson.M{}
	if d, err := f.MetadataFilter.ToBsonM(); err == nil {
		o.MetadataFilter = d
	} else {
		return o, err
	}

	// handling parameter
	o.HandlingParameter = f.HandlingParameter

	return o, nil
}
