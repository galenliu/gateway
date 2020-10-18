package controllers

import (
	"fmt"
	"gateway/models"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
)

var log *zap.Logger

type ThingsController struct {
	Container    *models.Things
	addonManager models.IAddonManager
}

func NewThingsController(manage models.IAddonManager, _log *zap.Logger) *ThingsController {
	if _log != nil {
		log = _log
	} else {
		log = zap.L()
	}
	return &ThingsController{
		Container:    models.NewThings(manage, _log),
		addonManager: manage,
	}
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
	log.Info("things controller handle create thing")
}

func (ts *ThingsController) HandleGetThing(c *gin.Context) {
	id := c.Param("thingId")
	log.Info("things handler:GetThing", zap.String("/:thingId", id), zap.String("method", "GET"))
	c.JSON(http.StatusOK, ts.Container.GetThing(id))
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
	var d []byte
	_, _ = c.Request.Body.Read(d)
	title := json.Get(d, "title").ToString()
	if title == "" {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	iconData := json.Get(d, "icon_data").ToString()
	if iconData != "" {
		thing.SetIcon(iconData, true)
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
