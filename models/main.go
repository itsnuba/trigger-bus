package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseId struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
}

type BaseCreatable struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type BaseModifable struct {
	ModifiedAt *time.Time `bson:"modifiedAt" json:"modifiedAt"`
}
