package requests

import (
	"time"

	"github.com/itsnuba/trigger-bus/commons"
	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerSchedulerForm struct {
	Metadata    *bson.M
	CronExpr    *string
	Active      *bool
	EndpointUrl *string
}

func (i TriggerSchedulerForm) ToCreateStruct() models.TriggerScheduler {
	o := models.MakeTriggerScheduler()

	o.Metadata = commons.NilOrValue(i.Metadata, bson.M{})
	o.CronExpr = commons.NilOrValue(i.CronExpr, "")
	o.Active = commons.NilOrValue(i.Active, true)
	o.EndpointUrl = commons.NilOrValue(i.EndpointUrl, "")

	return o
}

func (i TriggerSchedulerForm) ToUpdateMap() bson.M {
	o := bson.M{
		"modifiedAt": time.Now(),
	}

	commons.SetMapIfNotNil(o, "metadata", i.Metadata)
	commons.SetMapIfNotNil(o, "cronExpr", i.CronExpr)
	commons.SetMapIfNotNil(o, "active", i.Active)
	commons.SetMapIfNotNil(o, "endpointUrl", i.EndpointUrl)

	return o
}

type PostTriggerSchedulerForm struct {
	Metadata    bson.M `json:"" binding:""`
	CronExpr    string `json:"" binding:"required,cronExprFormat"`
	Active      *bool  `json:"" binding:""`
	EndpointUrl string `json:"" binding:"required,url"`
}

func (f PostTriggerSchedulerForm) ToForm() TriggerSchedulerForm {
	return TriggerSchedulerForm{
		Metadata:    &f.Metadata,
		CronExpr:    &f.CronExpr,
		Active:      f.Active,
		EndpointUrl: &f.EndpointUrl,
	}
}

type PutTriggerSchedulerForm struct {
	Metadata    *bson.M `json:"" binding:""`
	CronExpr    *string `json:"" binding:"omitempty,gt=0,cronExprFormat"`
	Active      *bool   `json:"" binding:""`
	EndpointUrl *string `json:"" binding:"omitempty,gt=0,url"`
}

func (f PutTriggerSchedulerForm) ToForm() TriggerSchedulerForm {
	return TriggerSchedulerForm(f)
}
