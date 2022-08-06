package queue

import "github.com/itsnuba/trigger-bus/models"

var EventChannel = make(chan models.Event, 10)
var SchedulerChannel = make(chan models.TriggerScheduler, 10)
