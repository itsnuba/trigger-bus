package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsnuba/trigger-bus/models"
	"github.com/itsnuba/trigger-bus/models/requests"
	"github.com/itsnuba/trigger-bus/models/responses"
	"github.com/itsnuba/trigger-bus/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const errStringDuplicateTriggerListener string = "similar listener already exist with id [%s]. \nif multiple listener is intended, use [?allow_duplicate=true]"

type getTriggerListenerListParam struct {
	Activity     string `form:"activity"`
	Microservice string `form:"microservice"`
}

type addTriggerListenerQueryParam struct {
	AllowDuplicate bool `form:"allow_duplicate"`
}

type editTriggerListenerParam struct {
	// path
	Id string `uri:"id" binding:"required,hexadecimal,len=24"`

	// query
	AllowDuplicate bool `form:"allow_duplicate"`
}

func getTriggerListener(cols *mongo.Collection, filter bson.M) *models.TriggerListener {
	var data models.TriggerListener
	if err := cols.FindOne(context.TODO(), filter).Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			panic(err)
		}
	}
	return &data
}

func GetTriggerListenerListHandler(c *gin.Context, cols *mongo.Collection) {
	var param getTriggerListenerListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err, "cannot parse query parameter"),
		)
		return
	}

	filters := bson.M{}
	if param.Activity != "" {
		filters["activity"] = param.Activity
	}
	if param.Microservice != "" {
		filters["metadata.microservice"] = param.Microservice
	}

	data := []models.TriggerListener{}
	if curr, err := cols.Find(context.TODO(), filters); err == nil {
		if err := curr.All(context.TODO(), &data); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}

func AddTriggerListenerHandler(c *gin.Context, cols *mongo.Collection) {
	var form requests.TriggerListenerAddForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err),
		)
		return
	}

	var qp addTriggerListenerQueryParam
	if err := c.ShouldBindQuery(&qp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err, "cannot parse query parameter"),
		)
		return
	}

	// prevent duplicate
	if !qp.AllowDuplicate {
		if ext := getTriggerListener(cols, bson.M{
			"activity":    form.Activity,
			"callbackUrl": form.CallbackUrl,
		}); ext != nil {
			c.AbortWithStatusJSON(http.StatusConflict,
				responses.MakeApiErrorResponse(fmt.Sprintf(errStringDuplicateTriggerListener, ext.Id.Hex())),
			)
			return
		}
	}

	data, err := form.ToTriggerListener()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			responses.MakeApiErrorResponseFromError(err),
		)
		return
	}
	data.Active = true // default awal regis, diisi true

	if _, err := cols.InsertOne(context.TODO(), data); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}

func EditTriggerListenerHandler(c *gin.Context, cols *mongo.Collection) {
	var param editTriggerListenerParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err, "cannot parse path parameter"),
		)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err, "cannot parse query parameter"),
		)
		return
	}

	id, err := primitive.ObjectIDFromHex(param.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			responses.MakeApiErrorResponse("cannot parse id", err.Error()),
		)
	}

	var form requests.TriggerListenerEditForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			validators.TranslateValidationError(err),
		)
		return
	}

	data := getTriggerListener(cols, bson.M{"_id": id})
	if data == nil {
		c.AbortWithStatusJSON(http.StatusNotFound,
			responses.MakeApiErrorResponse("not found"),
		)
		return
	}

	form.ApplyToTriggerListener(data)

	// cek duplikat
	if !param.AllowDuplicate {
		if dup := getTriggerListener(cols, bson.M{
			"_id":         bson.M{"$ne": id},
			"activity":    data.Activity,
			"callbackUrl": data.CallbackUrl,
		}); dup != nil {
			fmt.Println(dup)
			c.AbortWithStatusJSON(http.StatusConflict,
				responses.MakeApiErrorResponse(fmt.Sprintf(errStringDuplicateTriggerListener, dup.Id.Hex())),
			)
			return
		}
	}

	if _, err := cols.ReplaceOne(context.TODO(), bson.M{"_id": id}, data); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}
