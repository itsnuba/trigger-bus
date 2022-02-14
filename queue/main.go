package queue

import "github.com/itsnuba/trigger-bus/models"

var EventChannel = make(chan models.Event, 10)
