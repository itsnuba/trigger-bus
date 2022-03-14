package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/models/requests"
	"github.com/itsnuba/trigger-bus/models/responses"
	"github.com/itsnuba/trigger-bus/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddEventHandler(c *gin.Context, eventCols *mongo.Collection, listenerCols *mongo.Collection, logCols *mongo.Collection) {
	var form requests.EventAddForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			if e.Field == "payload" {
				c.AbortWithStatusJSON(http.StatusBadRequest,
					responses.MakeApiErrorResponse("payload must be list of object"),
				)
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err),
		)
		return
	}

	event, err := form.ToEvent()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			responses.MakeApiErrorResponseFromError(err),
		)
		return
	}

	if _, err := eventCols.InsertOne(context.TODO(), event); err != nil {
		panic(err)
	}

	// get listener
	listeners, _ := findListener(event, listenerCols)

	tls := []models.TriggerLog{}
	for _, l := range listeners {
		tl, err := publishEvent(l, event, logCols)
		if err != nil {
			panic(err)
		}
		tls = append(tls, *tl)
	}

	res := responses.AddEventResult{
		Event:    event,
		Handlers: tls,
	}

	c.JSON(http.StatusOK, res)
}

func findListener(event models.Event, listenerCols *mongo.Collection) ([]models.TriggerListener, error) {
	res := []models.TriggerListener{}

	// activity to filter
	activityPart := strings.Split(event.Activity, ";")
	activityFilter := []bson.M{}
	for i := range activityPart {
		activityFilter = append(activityFilter, bson.M{
			"activity": strings.Join(activityPart[0:i+1], ";"),
		})
	}

	// metadata filter
	metadataFilters := []bson.M{}
	for k, v := range event.Metadata {
		mk := "metadataFilter." + k
		mf := bson.M{
			"$or": []bson.M{
				{
					mk + ".0": bson.M{
						"$exists": false,
					},
				},
				{
					mk: bson.M{
						"$all": bson.A{v},
					},
				},
			},
		}
		metadataFilters = append(metadataFilters, mf)
	}

	curr, err := listenerCols.Find(context.TODO(), bson.M{
		"$and": append([]bson.M{
			{
				"$or": activityFilter,
			},
		}, metadataFilters...),
	})
	if err != nil {
		return res, err
	}

	if err := curr.All(context.TODO(), &res); err != nil {
		return res, err
	}

	println(len(res))

	return res, nil
}

func publishEvent(
	listener models.TriggerListener, event models.Event,
	logCols *mongo.Collection,
) (*models.TriggerLog, error) {
	// log
	log := listener.ToTriggerLog()
	log.EventId = event.Id

	if _, err := logCols.InsertOne(context.TODO(), log); err != nil {
		return nil, err
	}

	// send to handler using subroutine
	go func() {
		handled := false
		handledMessage := ""

		body, _ := json.Marshal(event.Payload)
		b := bytes.NewBuffer(body)
		if resp, err := http.Post(listener.CallbackUrl, "application/json", b); err != nil {
			handledMessage = err.Error()
		} else {
			defer resp.Body.Close()

			bb, _ := io.ReadAll(resp.Body)
			handled = resp.StatusCode == http.StatusOK
			handledMessage = fmt.Sprintf("code: %d\nbody:\n%s", resp.StatusCode, string(bb))
		}

		// update log
		if _, err := logCols.UpdateByID(context.TODO(), log.Id, bson.M{
			"$set": bson.M{
				"modifiedAt":     time.Now(),
				"handled":        handled,
				"handledMessage": handledMessage,
			},
		}); err != nil {
			fmt.Println(err)
		}
	}()

	return &log, nil
}
