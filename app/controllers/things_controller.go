package controllers

import (
	"fmt"
	"gateway/app/models"
	"gateway/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

func (tc *ThingsController) HandleWebsocket(c *gin.Context, thingId string) {

	tc.wsHandler(c.Writer, c.Request, thingId)

}

//GET  /things
func (tc *ThingsController) HandleGetThings(c *gin.Context) {
	if c.IsWebsocket() {
		tc.HandleWebsocket(c, "")
		return
	}
	log.Info("things controller handle get things")
	c.JSON(http.StatusOK, tc.Container.GetThings())
}

// POST /things
func (tc *ThingsController) HandleCreateThing(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}
	var thing models.Thing
	err = json.Unmarshal(data, &thing)
	if err != nil || thing.ID == "" {
		c.String(http.StatusBadRequest, "bad request")
	}
	if thing.SelectedCapability == "" {
		thing.SelectedCapability = "Custom"
	}
	err = tc.Container.CreateThing(thing)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("create thing(%s) err: %v", thing.ID, err.Error()))
	}
	c.String(http.StatusOK, fmt.Sprintf("create thing(%s) succeed", thing.ID))
}

func (tc *ThingsController) HandleGetThing(c *gin.Context) {

	thingId := c.Param("thing_id")
	if thingId == "" {
		c.String(http.StatusBadRequest, "thing id invalid")
		return
	}
	if c.IsWebsocket() {
		tc.HandleWebsocket(c, thingId)
		return
	}
	t, err := models.GetThing(thingId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, t)
}

//patch /
func (tc *ThingsController) HandlePatchThings(c *gin.Context) {
	log.Info("things controller handle patch things")

}

//patch /:thingId
func (tc *ThingsController) HandlePatchThing(c *gin.Context) {
	log.Info("things controller handle patch thing")

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

	v, err := tc.Container.SetThingProperty(thingId, propName, value)
	if err != nil {
		log.Error("Error setting property")
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, struct {
		PropertyName interface{} `json:"propertyName"`
	}{PropertyName: v})
	return
}

func (tc *ThingsController) HandleGetProperty(c *gin.Context) {
	thingId := c.Param("thing_id")
	propName := c.Param("property_name")
	value := tc.Container.GetThingProperty(thingId, propName)
	data := map[string]interface{}{propName: value}
	c.JSON(http.StatusOK, data)
}

func (tc *ThingsController) HandleGetProperties(c *gin.Context) {
	thingId := c.Param("thing_id")
	//thing := tc.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	//for propName, _ := range thing.Properties {
	//	props[propName] = tc.Container.Manager.GetProperty(thingId, propName)
	//}
	log.Info("things handler:GetProperties", zap.String("thingId", thingId), zap.String("method", "PUT"))
	c.JSON(http.StatusOK, props)
}

func (tc *ThingsController) HandleSetThing(c *gin.Context) {
	thingId := c.Param("thingId")
	thing, err := models.GetThing(thingId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	req, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Read body err")
		return
	}
	var rThing models.Thing
	err = json.Unmarshal(req, &rThing)
	if err != nil {
		c.String(http.StatusBadRequest, "Unmarshal err")
		return
	}
	if rThing.Title != "" {
		thing.SetTitle(rThing.Title)
	}

}

func (tc *ThingsController) HandleDeleteThing(c *gin.Context) {
	thingId := c.Param("thing_id")
	err := tc.Container.RemoveThing(thingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to remove thing thingId: %v ,err: %v", thingId, err))
		return
	}
	log.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
	c.Status(http.StatusNoContent)
}

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (tc *ThingsController) wsHandler(w http.ResponseWriter, r *http.Request, thingId string) {
	//conn, err := wsUpgrade.Upgrade(w, r, nil)
	//if err != nil {
	//}
	//thing,err := tc.Container.GetThing(thingId);
	//if err != nil {
	//
	//}

}

type ThingContext struct {
	conn    *websocket.Conn
	thingId string
}
