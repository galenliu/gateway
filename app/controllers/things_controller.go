package controllers

import (
	"fmt"
	"gateway/app/models"
	"gateway/util/logger"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

var log *zap.Logger

type ThingsController struct {
	Container *models.Things
}

func NewThingsController() *ThingsController {
	tc := &ThingsController{}
	tc.Container = models.NewThings()
	log = logger.GetLog()
	return tc
}

func (ts *ThingsController) HandleGetThings(c *gin.Context) {
	log.Info("things controller handle get things")
	c.JSON(http.StatusOK, ts.Container.GetThings())
}

//patch /
func (ts *ThingsController) HandlePatchThings(c *gin.Context) {
	log.Info("things controller handle patch things")

}

//patch /:thingId
func (ts *ThingsController) HandlePatchThing(c *gin.Context) {
	log.Info("things controller handle patch thing")

}

func (ts *ThingsController) HandleCreateThing(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("read request body err ")
		c.String(http.StatusBadRequest, "bad request")
	}
	var thing models.Thing
	err = json.Unmarshal(data, &thing)
	if err != nil {
		log.Error("read request body err ")
		c.String(http.StatusBadRequest, "bad request")
	}
	err = ts.Container.CreateThing(thing)
	if err != nil {
		log.Error(fmt.Sprintf("create thing err: %v", err))
	}
	c.Status(http.StatusOK)
}

func (ts *ThingsController) HandleGetThing(c *gin.Context) {
	id := c.Param("thingId")
	t  :=ts.Container.GetThing(id);if t==nil{
		c.JSON(http.StatusBadRequest, fmt.Sprintf("have not thing Id: %v",id))
	}
	log.Info("things handler:GetThing", zap.String("/:thingId", id), zap.String("method", "GET"))
	c.JSON(http.StatusOK,t)
}

//put property
func (ts *ThingsController) HandleSetProperty(c *gin.Context) {

	thingId := c.Param("thingId")
	propName := c.Param("propertyName")
	log.Info("things handler:SetProperty", zap.String("propName", propName), zap.String("method", "PUT"))

	var data []byte
	_, err := c.Request.Body.Read(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
	}
	valueInterface := json.Get(data, propName).GetInterface()
	ts.Container.SetThingProperty(thingId, propName, valueInterface)

}

func (ts *ThingsController) HandleGetProperty(c *gin.Context) {
	thingId := c.Param("thingId")
	propName := c.Param("propertyName")
	value := ts.Container.GetThingProperty(thingId, propName)
	data := map[string]interface{}{propName: value}
	c.JSON(http.StatusOK, data)
}

func (ts *ThingsController) HandleGetProperties(c *gin.Context) {
	thingId := c.Param("thingId")
	//thing := ts.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	//for propName, _ := range thing.Properties {
	//	props[propName] = ts.Container.Manager.GetProperty(thingId, propName)
	//}
	log.Info("things handler:GetProperties", zap.String("thingId", thingId), zap.String("method", "PUT"))
	c.JSON(http.StatusOK, props)
}

func (ts *ThingsController) HandleSetThing(c *gin.Context) {
	thingId := c.Param("thingId")
	thing := ts.Container.GetThing(thingId)

	req, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(fmt.Sprintf("set thing err, request read err :%v", err))
		c.JSON(http.StatusBadRequest, "")
	}
	title := json.Get(req, "title").ToString()
	if title == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	err = ts.Container.SetThing(*thing)
	if err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}
}


func (ts *ThingsController) HandleDeleteThing(c *gin.Context) {
	thingId := c.Param("thingId")
	err := ts.Container.RemoveThing(thingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to remove thing thingId: %v ,err: %v", thingId, err))
		return
	}
	log.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
	c.Status(http.StatusNoContent)
}
