package queue

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/itsnuba/trigger-bus/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DoSchedulerEvent(schCols *mongo.Collection, logCols *mongo.Collection) {
	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()
	s.StartAsync()

	// regis existing active scheduler
	go regisScheduler(schCols)

	for sch := range SchedulerChannel {
		tag := sch.Id.Hex()
		if sch.Active {
			// regis
			s.Cron(sch.CronExpr).Tag(tag).Do(func(l *mongo.Collection, s models.TriggerScheduler) {
				doScheduler(l, s)
			}, logCols, sch)
		} else {
			// unregis
			s.RemoveByTag(tag)
		}
		fmt.Println(len(s.Jobs()))
	}
}

func regisScheduler(schCols *mongo.Collection) {
	var data []models.TriggerScheduler
	if curr, err := schCols.Find(context.TODO(), bson.M{"active": true}); err != nil {
		return
	} else if err = curr.All(context.TODO(), &data); err != nil {
		return
	}

	for _, d := range data {
		SchedulerChannel <- d
	}
}

func doScheduler(logCols *mongo.Collection, sch models.TriggerScheduler) {
	if err := func() (err error) {
		// log
		log := sch.ToTriggerLog()
		if _, err = logCols.InsertOne(context.TODO(), log); err != nil {
			return
		}

		handled := false
		handledMessage := ""

		if resp, err := http.Get(sch.EndpointUrl); err != nil {
			handledMessage = err.Error()
		} else {
			defer resp.Body.Close()

			bb, _ := io.ReadAll(resp.Body)
			handled = resp.StatusCode == http.StatusOK
			handledMessage = fmt.Sprintf("code: %d\nbody:\n%s", resp.StatusCode, string(bb))
		}

		// update log
		if _, err = logCols.UpdateByID(context.TODO(), log.Id, bson.M{
			"$set": bson.M{
				"modifiedAt":     time.Now(),
				"handled":        handled,
				"handledMessage": handledMessage,
			},
		}); err != nil {
			return
		}
		return
	}(); err != nil {
		fmt.Printf("failed to do scheduler [%s]. cause:\n", sch.Id)
		fmt.Println(err)
	}
}
