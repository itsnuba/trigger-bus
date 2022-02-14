package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient      *mongo.Client
	dbCollections struct {
		events           *mongo.Collection
		triggerListeners *mongo.Collection
		triggerLogs      *mongo.Collection
	}
)

func setupDB() {
	opts := options.Client()
	opts.ApplyURI(config.MongoUri)
	opts.SetConnectTimeout(time.Second * 30)
	if config.Debug {
		opts.SetMonitor(&event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		})
	}

	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal(err)
	}

	// replica set gabisa di ping
	// https://www.mongodb.com/community/forums/t/go-driver-cant-connect-to-an-uninitialized-replica-set/104469/2
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	dbClient = client
	dbCollections.events = client.Database(config.MongoDB).Collection(config.MongoColEvent)
	dbCollections.triggerListeners = client.Database(config.MongoDB).Collection(config.MongoColTriggerListener)
	dbCollections.triggerLogs = client.Database(config.MongoDB).Collection(config.MongoColTriggerLog)
}
