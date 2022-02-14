package queue

import (
	"context"
	"fmt"

	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DoPublishEvent(listenerCols *mongo.Collection, logCols *mongo.Collection) {
	for event := range EventChannel {
		fmt.Println(len(EventChannel))

		err := publishEvent(event, listenerCols, logCols)
		if err != nil {
			fmt.Println("[EVENT-PUBLISH] failed")
			fmt.Println(err)
		}
	}

	fmt.Println("[EVENT-PUBLISH] channel broken")
}

func publishEvent(event models.Event, listenerCols *mongo.Collection, logCols *mongo.Collection) error {
	listeners, err := findListener(event, listenerCols)
	if err != nil {
		return err
	}

	fmt.Printf("found %d listener\n", len(listeners))
	for _, l := range listeners {
		fmt.Println(l.CallbackUrl)
	}

	return nil
}

func findListener(event models.Event, listenerCols *mongo.Collection) ([]models.TriggerListener, error) {
	accountIds := []int{}
	if v, ok := event.Metadata["accountId"]; ok {
		if id, ok := v.(float64); ok {
			accountIds = append(accountIds, int(id))
		}
	}

	curr, err := listenerCols.Find(context.TODO(), bson.M{
		"activity": event.Activity,
		"$or": []bson.M{
			{
				"metadataFilter.accountIds.0": bson.M{
					"$exists": false,
				},
			},
			{
				"metadataFilter.accountIds": bson.M{
					"$all": accountIds,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	res := []models.TriggerListener{}
	if err := curr.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	return res, nil
}
