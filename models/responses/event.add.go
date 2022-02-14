package responses

import "github.com/itsnuba/trigger-bus/models"

type AddEventResult struct {
	models.Event

	Handlers []models.TriggerLog `json:"handlers"`
}
