package main

import (
	"fmt"
	"time"

	"github.com/itsnuba/trigger-bus/configs"
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/utils/operations"
)

// config
var config *configs.Config

func main() {
	config = configs.Load()
	setupDB()

	consoleHelper(
		Options{
			"trigger_logs",
		},
		func() bool {
			consoleHelper(
				Options{
					"delete",
				},
				func() bool {
					triggerTypeS := readInputAsString("type [", string(models.TLETEvent), ", ", string(models.TLETScheduler), "]: ")
					triggerType := models.TrigglerLogTriggerType(triggerTypeS)
					if triggerType != models.TLETEvent && triggerType != models.TLETScheduler {
						fmt.Println("type not found")
						return false
					}

					var err error
					var since, until time.Time

					sinceS := readInputAsString("since [empty, yyyy-mm-dd]: ")
					if sinceS != "" {
						if since, err = time.Parse("2006-01-02", sinceS); err != nil {
							fmt.Printf("cannot parse since, %s\n", err)
							return false
						}
					}
					untilS := readInputAsString("until [empty, yyyy-mm-dd]: ")
					if untilS != "" {
						if until, err = time.Parse("2006-01-02", untilS); err != nil {
							fmt.Printf("cannot parse until, %s\n", err)
							return false
						}
					}

					if err := operations.DeleteLogConfirmation(dbCollections.triggerLogs, triggerType, since, until); err != nil {
						fmt.Println(err)
						return false
					}

					if confirm := readInputAsString("type [confirm] to delete logs: "); confirm != "confirm" {
						fmt.Println("aborting delete logs")
						return false
					}

					operations.DeleteLog(dbCollections.triggerLogs, triggerType, since, until)

					return true
				},
			)

			return true
		},
	)

	readInputAsString("press [enter] to exit")
}
