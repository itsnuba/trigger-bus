package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/handlers"
	m "github.com/itsnuba/trigger-bus/handlers/middlewares"
)

func setupRouter(r *gin.Engine) {
	h := handlers.MakeHandler(
		dbClient,
		dbCollections.events,
		dbCollections.triggerListeners,
		dbCollections.triggerScheduler,
		dbCollections.triggerLogs,
	)

	// CORS
	r.Use(cors.Default())

	// health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// event
	r.POST("/events", func(c *gin.Context) {
		handlers.AddEventHandler(c, dbCollections.events, dbCollections.triggerListeners, dbCollections.triggerLogs)
	})

	// trigger listener
	{
		g := r.Group("/trigger_listeners")
		g.GET("", func(c *gin.Context) {
			handlers.GetTriggerListenerListHandler(c, dbCollections.triggerListeners)
		})
		g.POST("", func(c *gin.Context) {
			handlers.AddTriggerListenerHandler(c, dbCollections.triggerListeners)
		})
		g.GET("/:id", func(c *gin.Context) {
			handlers.GetTriggerListenerByIdHandler(c, dbCollections.triggerListeners)
		})
		g.PUT("/:id", func(c *gin.Context) {
			handlers.EditTriggerListenerHandler(c, dbCollections.triggerListeners)
		})
		g.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteTriggerListenerHandler(c, dbCollections.triggerListeners)
		})
	}

	// trigger handling
	{
		g := r.Group("/handle")
		g.PUT("", func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		})
	}

	// trigger scheduler
	{
		g := r.Group("/trigger_schedulers")
		g.GET("", h.GetTriggerScheduler)
		g.POST("", h.PostTriggerScheduler)

		gId := g.Group("/:id", m.GetResourceId)
		gId.GET("", h.GetTriggerSchedulerById)
		gId.PUT("", h.PutTriggerSchedulerById)
		gId.DELETE("", h.DeleteTriggerListenerById)
	}

}
