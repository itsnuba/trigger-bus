package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/handlers/helpers"
	"github.com/itsnuba/trigger-bus/handlers/middlewares"
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/models/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type triggerSchedulerParam struct {
	AllowDuplicate bool `form:"allow_duplicate"`
	Reload         bool `form:"reload"`
}

func (h *Handlers) GetTriggerScheduler(c *gin.Context) {
	var param struct {
		Metadata string `form:"metadata"`
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		helpers.HandleParsingError(c, err, "invalid query parameter")
		return
	}

	filters := bson.M{}

	data, err := getTriggerScheduler(context.TODO(), h.triggerSchedulersCol, filters)
	if err != nil {
		helpers.HandleError(c, err)
		return
	}

	// reload?
	go func() {
		var qp triggerSchedulerParam
		c.ShouldBindQuery(&qp)
		if qp.Reload {
			for _, d := range data {
				helpers.SendSchedulerToSchedulerChannel(d)
			}
		}
	}()

	if data == nil {
		data = []models.TriggerScheduler{}
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handlers) GetTriggerSchedulerById(c *gin.Context) {
	id := c.MustGet(middlewares.ResourceIdS).(primitive.ObjectID)

	data, err := getTriggerSchedulerById(context.TODO(), h.triggerSchedulersCol, id)
	if err != nil {
		helpers.HandleError(c, err)
		return
	}

	// reload?
	go func() {
		var qp triggerSchedulerParam
		c.ShouldBindQuery(&qp)
		if qp.Reload {
			helpers.SendSchedulerToSchedulerChannel(data)
		}
	}()

	c.JSON(http.StatusOK, data)
}

func (h *Handlers) PostTriggerScheduler(c *gin.Context) {
	var req requests.PostTriggerSchedulerForm
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleParsingError(c, err, "invalid body")
		return
	}

	var qp triggerSchedulerParam
	c.ShouldBindQuery(&qp)

	form := req.ToForm()

	data, err := createTriggerScheduler(context.TODO(), h.triggerSchedulersCol, form, qp.AllowDuplicate)
	if err != nil {
		if errors.Is(err, helpers.ErrDuplicate) {
			err = fmt.Errorf("%w. \nif multiple scheduler is intended, use [?allow_duplicate=true]", err)
		}
		helpers.HandleError(c, err)
		return
	}

	// regis scheduler
	go helpers.SendSchedulerToSchedulerChannel(data)

	c.JSON(http.StatusOK, data)
}

func (h *Handlers) PutTriggerSchedulerById(c *gin.Context) {
	id := c.MustGet(middlewares.ResourceIdS).(primitive.ObjectID)

	var req requests.PutTriggerSchedulerForm
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleParsingError(c, err, "invalid body")
		return
	}

	var qp triggerSchedulerParam
	c.ShouldBindQuery(&qp)

	form := req.ToForm()

	data, err := editTriggerSchedulerById(context.TODO(), h.triggerSchedulersCol, id, form, qp.AllowDuplicate)
	if err != nil {
		if errors.Is(err, helpers.ErrDuplicate) {
			err = fmt.Errorf("%w. \nif multiple scheduler is intended, use [?allow_duplicate=true]", err)
		}
		helpers.HandleError(c, err)
		return
	}

	// reload scheduler
	go helpers.SendSchedulerToSchedulerChannel(data)

	c.JSON(http.StatusOK, data)
}

func (h *Handlers) DeleteTriggerListenerById(c *gin.Context) {
	id := c.MustGet(middlewares.ResourceIdS).(primitive.ObjectID)

	data, err := deleteTriggerScheduler(context.TODO(), h.triggerSchedulersCol, id)
	if err != nil {
		helpers.HandleError(c, err)
		return
	}

	data.Active = false

	// regis scheduler
	go helpers.SendSchedulerToSchedulerChannel(data)

	c.JSON(http.StatusOK, data)
}
