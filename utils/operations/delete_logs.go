package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteLogConfirmation(
	cols *mongo.Collection,
	triggerType models.TrigglerLogTriggerType,
	since time.Time, until time.Time,
) error {
	pipe := []bson.M{}

	// filter
	filters := bson.M{}
	if !since.IsZero() {
		filters["$gte"] = bson.A{
			"$createdAt", since,
		}
	}
	if !until.IsZero() {
		filters["$lte"] = bson.A{
			"$createdAt", until,
		}
	}

	// filter
	pipe = append(pipe, bson.M{
		"$match": bson.M{
			"triggerType": triggerType,
		},
	})

	// group
	pipe = append(pipe, bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"$cond": bson.A{
					filters,
					"deleted",
					"remaining",
				},
			},
			"tot": bson.M{
				"$sum": 1,
			},
		},
	})

	data := []struct {
		Title string `bson:"_id"`
		Total int    `bson:"tot"`
	}{}
	if curr, err := cols.Aggregate(context.TODO(), pipe); err == nil {
		if err := curr.All(context.TODO(), &data); err != nil {
			return err
		}
	} else {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("logs not found")
	}

	var deleted, remaining int
	for _, d := range data {
		if d.Title == "deleted" {
			deleted = d.Total
		} else if d.Title == "remaining" {
			remaining = d.Total
		}
	}

	fmt.Println("delete confirmation: ")
	fmt.Printf("%d of %d logs ", deleted, deleted+remaining)
	if !since.IsZero() {
		fmt.Printf("since %s ", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		fmt.Printf("until %s ", until.Format("2006-01-02"))
	}
	fmt.Println("will be deleted")

	return nil
}

func DeleteLog(
	cols *mongo.Collection,
	triggerType models.TrigglerLogTriggerType,
	since time.Time, until time.Time,
) (bool, error) {
	// filter
	filters := bson.M{
		"triggerType": triggerType,
	}
	// since until
	if !since.IsZero() || !until.IsZero() {
		createdFilters := bson.M{}
		if !since.IsZero() {
			createdFilters["$gte"] = since
		}
		if !until.IsZero() {
			createdFilters["$lte"] = until
		}
		filters["createdAt"] = createdFilters
	}

	if res, err := cols.DeleteMany(context.TODO(), filters); err != nil {
		return false, err
	} else {
		fmt.Printf("successfully delete %d logs\n", res.DeletedCount)
		return true, nil
	}
}
