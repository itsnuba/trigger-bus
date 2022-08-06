package helpers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/models/responses"
	"github.com/itsnuba/trigger-bus/validators"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDuplicate = errors.New("duplicate")

func HandleError(c *gin.Context, err error) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound,
			responses.MakeApiErrorResponse("not found"))
	} else {
		c.JSON(http.StatusUnprocessableEntity,
			responses.MakeApiErrorResponseFromError(err))
	}

	c.Abort()
}

func HandleParsingError(c *gin.Context, err error, additionalMessage ...string) {
	c.JSON(http.StatusBadRequest,
		validators.TranslateValidationError(err, additionalMessage...))

	c.Abort()
}
