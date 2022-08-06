package helpers

import (
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/queue"
)

func SendEventToEventChannel(event models.Event) {
	queue.EventChannel <- event
}

func SendSchedulerToSchedulerChannel(sch models.TriggerScheduler) {
	queue.SchedulerChannel <- sch
}
