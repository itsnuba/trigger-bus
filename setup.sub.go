package main

import (
	"github.com/itsnuba/trigger-bus/queue"
)

func regisGoSub() {
	// event publisher
	go queue.DoPublishEvent(dbCollections.triggerListeners, dbCollections.triggerLogs)
}
