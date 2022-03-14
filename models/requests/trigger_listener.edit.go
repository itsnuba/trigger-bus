package requests

import (
	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerListenerEditApplayer interface {
	ApplyToTriggerListener(o *models.TriggerListener) error
}

type TriggerListenerEditForm struct {
	Activity          *string                        `json:"activity" binding:"omitempty,activityFormat"`
	CallbackUrl       *string                        `json:"callbackUrl" binding:"omitempty,url"`
	Active            *bool                          `json:"active"`
	MetadataFilter    *TriggerListenerMetadataFilter `json:"metadataFilter" binding:"omitempty,metadataFilterFormat"`
	HandlingParameter *bson.M                        `bson:"handlingParameter"`
}

func (f TriggerListenerEditForm) ApplyToTriggerListener(o *models.TriggerListener) error {
	if f.Active != nil {
		o.Active = *f.Active
	}
	if f.Activity != nil {
		o.Activity = *f.Activity
	}
	if f.CallbackUrl != nil {
		o.CallbackUrl = *f.CallbackUrl
	}
	if f.HandlingParameter != nil {
		o.HandlingParameter = *f.HandlingParameter
	}
	if f.MetadataFilter != nil {
		o.MetadataFilter = *f.MetadataFilter
	}

	return nil
}
