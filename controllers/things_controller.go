package controllers

import (
	"gateway/plugin"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"net/http"
)
import "gateway/models"

type ThingsController struct {
	Container *models.Things
}

func NewThingsController(manage *plugin.AddonsManager) *ThingsController {
	return &ThingsController{
		Container: models.NewThings(manage),
	}
}

func (ts *ThingsController) HandleGetThings(c *gin.Context) {
	c.JSON(http.StatusOK, ts.Container.GetThings())
}

//patch /
func (ts *ThingsController) HandlePatchThings(c *gin.Context) {

}

//patch /:thingId
func (ts *ThingsController) HandlePatchThing(c *gin.Context) {

}

func (ts *ThingsController) HandleCreateThing(c *gin.Context) {

}

func (ts *ThingsController) HandleGetThing(c *gin.Context) {
	id := c.Param("thingId")
	c.JSON(http.StatusOK, ts.Container.GetThing(id))
}

//put property
func (ts *ThingsController) HandleSetProperty(c *gin.Context) {
	thingId := c.Param("thingId")
	propName := c.Param("propertyName")

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
	thing := ts.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	for propName, _ := range thing.Properties {
		props[propName] = ts.Container.AddonManager.GetProperty(thingId, propName)
	}
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

}
