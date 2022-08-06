package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/itsnuba/trigger-bus/handlers/helpers"
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/models/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getTriggerScheduler(
	ctx context.Context,
	cols *mongo.Collection,
	filters bson.M,
) (res []models.TriggerScheduler, err error) {
	var curr *mongo.Cursor
	if curr, err = cols.Find(context.TODO(), filters); err != nil {
		return
	}

	if err = curr.All(context.TODO(), &res); err != nil {
		return
	}

	return
}

func getTriggerSchedulerById(
	ctx context.Context,
	cols *mongo.Collection,
	id primitive.ObjectID,
) (res models.TriggerScheduler, err error) {
	if err = cols.FindOne(ctx, bson.M{"_id": id}).Decode(&res); err != nil {
		return
	}

	return
}

func checkDuplicateTriggerScheduler(
	ctx context.Context,
	cols *mongo.Collection,
	form requests.TriggerSchedulerForm,
	id primitive.ObjectID,
) (err error) {
	var match models.TriggerScheduler
	if err = cols.FindOne(ctx, bson.M{
		"_id":         bson.M{"$ne": id},
		"endpointUrl": form.EndpointUrl,
		// "cronExpr":    form.CronExpr,
	}).Decode(&match); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
	} else {
		err = fmt.Errorf("%w: similar scheduler already exist with id [%s]", helpers.ErrDuplicate, match.Id.Hex())
	}

	return
}

func createTriggerScheduler(
	ctx context.Context,
	cols *mongo.Collection,
	form requests.TriggerSchedulerForm,
	allowDuplicate bool,
) (res models.TriggerScheduler, err error) {
	// prevent duplicate
	if !allowDuplicate {
		if err = checkDuplicateTriggerScheduler(ctx, cols, form, primitive.NilObjectID); err != nil {
			return
		}
	}

	res = form.ToCreateStruct()
	if _, err = cols.InsertOne(ctx, res); err != nil {
		return
	}

	return
}

func editTriggerSchedulerById(
	ctx context.Context,
	cols *mongo.Collection,
	id primitive.ObjectID, form requests.TriggerSchedulerForm,
	allowDuplicate bool,
) (res models.TriggerScheduler, err error) {
	// prevent duplicate
	if !allowDuplicate {
		if err = checkDuplicateTriggerScheduler(ctx, cols, form, id); err != nil {
			return
		}
	}

	data := form.ToUpdateMap()
	if _, err = cols.UpdateByID(ctx, id, bson.M{"$set": data}); err != nil {
		return
	}

	if res, err = getTriggerSchedulerById(ctx, cols, id); err != nil {
		return
	}

	return
}

func deleteTriggerScheduler(
	ctx context.Context,
	cols *mongo.Collection,
	id primitive.ObjectID,
) (res models.TriggerScheduler, err error) {
	if res, err = getTriggerSchedulerById(ctx, cols, id); err != nil {
		return
	}

	if _, err = cols.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return
	}

	return
}
