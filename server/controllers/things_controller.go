// @Title  things_controllers
// @Description  app router
// @Author  liuguilin
// @Update  liuguilin

package controllers

import (
	"fmt"
	"gateway/pkg/log"
	"gateway/plugin"
	"gateway/server/models"
	thing2 "gateway/server/models/thing"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type ThingsController struct {
	Container *models.Things
}

func NewThingsControllerFunc() *ThingsController {
	tc := &ThingsController{}
	tc.Container = models.NewThings()
	return tc
}

// POST /things
func (tc *ThingsController) handleCreateThing(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	log.Debug("Post /thing,Body: \t\n %s", data)

	id := json.Get(data, "id").ToString()
	if len(id) < 1 {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	t := tc.Container.GetThing(id)
	if t != nil {
		c.String(http.StatusBadRequest, "thing already added")
		return
	}

	des, err1 := tc.Container.CreateThing(id, data)
	if err1 != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("create thing(%s) err: %v", id, err.Error()))
		return
	}
	c.String(http.StatusCreated, des)
}

// DELETE /things/:thingId
func (tc *ThingsController) handleDeleteThing(c *gin.Context) {
	thingId := c.Param("thingId")
	_ = plugin.RemoveDevice(thingId)
	err := tc.Container.RemoveThing(thingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to remove thing thingId: %v ,err: %v", thingId, err))
		return
	}
	log.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
	c.Status(http.StatusNoContent)
}

//GET /things/:thingId
func (tc *ThingsController) handleGetThing(c *gin.Context) {
	if c.IsWebsocket() {
		handleWebsocket(c, tc.Container)
		return
	}

	thingId := c.Param("thingId")
	if thingId == "" {
		c.String(http.StatusBadRequest, "thing id invalid")
		return
	}
	t := tc.Container.GetThing(thingId)
	if t == nil {
		c.String(http.StatusBadRequest, "thing not found")
		return
	}
	c.JSON(http.StatusOK, t)
}

//GET /things
func (tc *ThingsController) handleGetThings(c *gin.Context) {
	if c.IsWebsocket() {
		handleWebsocket(c, tc.Container)
		return
	}
	ts := tc.Container.GetListThings()
	c.JSON(http.StatusOK, ts)
}

//patch things
func (tc *ThingsController) handlePatchThings(c *gin.Context) {
	log.Info("container controller handle patch container")

}

//PATCH /things/:thingId
func (tc *ThingsController) handlePatchThing(c *gin.Context) {
	log.Info("container controller handle patch thing")

}

//PUT /:thing/properties/:propertyName
func (tc *ThingsController) handleSetProperty(c *gin.Context) {

	thingId := c.Param("thingId")
	propName := c.Param("propertyName")

	if thingId == "" || propName == "" {
		c.String(http.StatusBadRequest, "invalid params")
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid params")
		return
	}
	value := json.Get(data, propName).GetInterface()

	prop, e := tc.Container.SetThingProperty(thingId, propName, value)
	if e != nil {
		log.Error("Failed set thing(%s) property:(%s) value:(%v),err:(%s)", thingId, propName, value, e.Error())
		c.JSON(http.StatusInternalServerError, struct {
			Value interface{} `json:"value"`
		}{Value: value})
		return
	}
	c.JSON(http.StatusOK, struct {
		Value interface{} `json:"value"`
	}{Value: prop.Value})
	return
}

func (tc *ThingsController) handleGetProperty(c *gin.Context) {
	thingId := c.Param("thing_id")
	propName := c.Param("property_name")
	property := tc.Container.FindThingProperty(thingId, propName)
	if property == nil {
		c.String(http.StatusBadRequest, "property no found")
		return
	}
	data := map[string]interface{}{propName: property.Value}
	c.JSON(http.StatusOK, data)
}

func (tc *ThingsController) handleGetProperties(c *gin.Context) {
	thingId := c.Param("thing_id")
	//thing := tc.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	//for propName, _ := range thing.Properties {
	//	props[propName] = tc.Container.Manager.GetProperty(thingId, propName)
	//}
	log.Info("container handler:GetProperties", zap.String("thingId", thingId), zap.String("method", "PUT"))
	c.JSON(http.StatusOK, props)
}

func (tc *ThingsController) handleSetThing(c *gin.Context) {
	thingId := c.Param("thingId")
	thing := tc.Container.GetThing(thingId)
	if thing == nil {
		c.String(http.StatusBadRequest, "thing not found")
		return
	}
	req, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Read body err")
		return
	}
	var rThing thing2.Thing
	err = json.Unmarshal(req, &rThing)
	if err != nil {
		c.String(http.StatusBadRequest, "Unmarshal err")
		return
	}
	if rThing.Title != "" {
		thing.SetTitle(rThing.Title)
	}

}
