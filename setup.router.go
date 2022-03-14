package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/handlers"
)

func setupRouter(r *gin.Engine) {

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

}
