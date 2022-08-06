package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/handlers/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ResourceIdS string = ":id" // type primitive.ObjectID

func GetResourceId(c *gin.Context) {
	var pp struct {
		Id string `uri:"id" binding:"required,hexadecimal,len=24"`
	}

	if err := c.ShouldBindUri(&pp); err != nil {
		helpers.HandleError(c, err)
		return
	}

	id, err := primitive.ObjectIDFromHex(pp.Id)
	if err != nil {
		helpers.HandleError(c, fmt.Errorf("cannot parse id")) // err belum ikutan
	}

	c.Set(ResourceIdS, id)
}
