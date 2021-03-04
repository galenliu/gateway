package controllers

import (
	"fmt"
	"gateway/pkg/log"
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
func (tc *ThingsController) HandleCreateThing(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	log.Debug("Post /thing,Body: \t\n %s", data)

	id := json.Get(data,"id").ToString()
	if len(id)<1 {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if tc.Container.HasThing(id) {
		c.String(http.StatusBadRequest, "thing already added")
		return
	}

	des,err1 := tc.Container.CreateThing(id,data)
	if err1 != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("create thing(%s) err: %v", err.Error()))
		return
	}
	c.String(http.StatusCreated, des)
}

func (tc *ThingsController) HandleDeleteThing(c *gin.Context) {
	thingId := c.Param("thingId")
	err := tc.Container.RemoveThing(thingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to remove thing thingId: %v ,err: %v", thingId, err))
		return
	}
	log.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
	c.Status(http.StatusNoContent)
}



func (tc *ThingsController) HandleGetThing(c *gin.Context) {
	if c.IsWebsocket() {
		HandleWebsocket(c, tc.Container)
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

func (tc *ThingsController) HandleGetThings(c *gin.Context) {
	if c.IsWebsocket() {
		HandleWebsocket(c, tc.Container)
		return
	}
	ts := tc.Container.GetThings()
	c.JSON(http.StatusOK, ts)
}

//patch /
func (tc *ThingsController) HandlePatchThings(c *gin.Context) {
	log.Info("container controller handle patch container")

}

//patch /:thingId
func (tc *ThingsController) HandlePatchThing(c *gin.Context) {
	log.Info("container controller handle patch thing")

}

func (tc *ThingsController) HandleSetProperty(c *gin.Context) {
	//PUT /:thing/properties/:propertyName
	thingId := c.Param("thingId")
	propName := c.Param("propertyName")

	if thingId == "" || propName == "" {
		c.String(http.StatusBadRequest, "invalid params")
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid value")
		return
	}
	value := json.Get(data, propName).GetInterface()

	tc.Container.SetThingProperty(thingId, propName, value)

	c.JSON(http.StatusOK, struct {
		Value interface{} `json:"value"`
	}{Value: value})
	return
}

func (tc *ThingsController) HandleGetProperty(c *gin.Context) {
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

func (tc *ThingsController) HandleGetProperties(c *gin.Context) {
	thingId := c.Param("thing_id")
	//thing := tc.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	//for propName, _ := range thing.Properties {
	//	props[propName] = tc.Container.Manager.GetProperty(thingId, propName)
	//}
	log.Info("container handler:GetProperties", zap.String("thingId", thingId), zap.String("method", "PUT"))
	c.JSON(http.StatusOK, props)
}

func (tc *ThingsController) HandleSetThing(c *gin.Context) {
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



