package requests

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type TriggerListenerMetadataFilter struct {
	AccountIds []int `json:"accountIds" binding:""`
}

func (v *TriggerListenerMetadataFilter) FromBsonM(i bson.M) error {
	if d, err := json.Marshal(i); err == nil {
		if err := json.Unmarshal(d, v); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (v TriggerListenerMetadataFilter) ToBsonM() (bson.M, error) {
	o := bson.M{}
	if d, err := json.Marshal(v); err == nil {
		if err := json.Unmarshal(d, &o); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return o, nil
}
