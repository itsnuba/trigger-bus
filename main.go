package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/itsnuba/trigger-bus/configs"
	"github.com/itsnuba/trigger-bus/validators"
)

// config
var config *configs.Config

func main() {
	config = configs.Load()
	setupDB()

	// setup gin custom validation & language
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.InitValidator(v)
	}

	// setup gin
	r := gin.New()
	// disable log for /ping
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{
			"/ping",
		},
	}))
	r.Use(gin.Recovery())

	setupRouter(r)

	// regis go subroutine
	regisGoSub()

	// run service
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
