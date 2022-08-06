package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handlers struct {
	db                   *mongo.Client
	eventsCol            *mongo.Collection
	triggerListenersCol  *mongo.Collection
	triggerSchedulersCol *mongo.Collection
	triggerLogsCol       *mongo.Collection
}

func MakeHandler(
	db *mongo.Client,
	eventsCol *mongo.Collection,
	triggerListenersCol *mongo.Collection,
	triggerSchedulersCol *mongo.Collection,
	triggerLogsCol *mongo.Collection,
) Handlers {
	return Handlers{
		db,
		eventsCol,
		triggerListenersCol,
		triggerSchedulersCol,
		triggerLogsCol,
	}
}
