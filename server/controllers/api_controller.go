package controllers

import (
	"addon"
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/plugin"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"strings"

	"sync"
)

type ApiController struct {
	thingId            string
	Container          *models.Things
	ws                 *websocket.Conn
	locker             *sync.Mutex
	done               chan struct{}
	subscriptionThings map[string]*thing.Thing
}

func NewApiController() *ThingsWebsocketHandler {
	controller := &ThingsWebsocketHandler{}
	controller.subscriptionThings = make(map[string]*thing.Thing, 10)
	controller.locker = new(sync.Mutex)
	controller.Container = models.NewThings()
	controller.done = make(chan struct{})
	return controller
}

func handleApiWebsocket(c *websocket.Conn) {

	if !c.Locals("websocket").(bool) {
		return
	}
	thingId := c.Locals("thingId").(string)

	controller := NewApiController()
	controller.ws = c

	_ = bus.Subscribe(util.PropertyChanged, controller.onPropertyChanged)
	_ = bus.Subscribe(util.CONNECTED, controller.onConnected)
	_ = bus.Subscribe(util.MODIFIED, controller.onModified)
	_ = bus.Subscribe(util.ThingRemoved, controller.onRemoved)

	for {
		select {
		case <-controller.done:
			controller.close()
			return
		default:
			_, data, readErr := controller.ws.ReadMessage()
			if readErr != nil {
				log.Info("websocket disconnected err: %s", readErr.Error())
				controller.done <- struct{}{}
				return
			}
			if data != nil {
				go controller.handleMessage(data)
			}

		}

	}

}

func (controller *ApiController) handleMessage(data []byte) {
	messageType := json.Get(data, "messageType").ToString()

	if messageType == "" {
		controller.sendMessage(struct {
			MessageType string      `json:"messageType"`
			Data        interface{} `json:"data"`
		}{
			MessageType: thing.Error,
			Data: struct {
				Code    int    `json:"code"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Code:    http.StatusBadRequest,
				Status:  "400 Bed Request",
				Message: fmt.Sprintf("messageType(%s) err", messageType),
			},
		})
		return
	}

	switch messageType {
	case models.GetThings:

		m := make(map[string]interface{})
		m["id"] = id
		m["messageType"] = models.ThingModified
		return
		return
	}

	switch messageType {
	case models.SetProperty:
		var propertyMap map[string]interface{}
		json.Get(data, "data").ToVal(&propertyMap)

		for propName, value := range propertyMap {
			p, setErr := plugin.SetProperty(device.GetID(), propName, value)

			m := make(map[string]interface{})
			m["messageType"] = thing.Error

			if setErr != nil {
				m["messageType"] = thing.Error
				m["data"] = map[string]interface{}{
					"code":    http.StatusBadRequest,
					"status":  "400 Bed Request",
					"message": setErr.Error(),
					"request": json.Get(data),
				}

			} else {

				m["messageType"] = thing.Error
				m["data"] = map[string]interface{}{
					"code":    http.StatusBadRequest,
					"status":  "400 Bed Request",
					"message": fmt.Sprintf("Unknown messageType:%s", messageType),
					"request": json.Get(data),
				}

				controller.sendMessage(struct {
					MessageType string      `json:"messageType"`
					Data        interface{} `json:"data"`
				}{
					MessageType: thing.Error,
					Data: struct {
						Code    int         `json:"code"`
						Status  string      `json:"status"`
						Message interface{} `json:"message"`
						Request interface{} `json:"request"`
					}{
						Code:   http.StatusOK,
						Status: "200 OK",
						Message: struct {
							Value interface{} `json:"value"`
						}{
							Value: json.Get(p, "value").GetInterface(),
						},
						Request: json.Get(data),
					},
				})

			}
		}
		return
	default:
		m := make(map[string]interface{})
		m["messageType"] = thing.Error
		m["data"] = map[string]interface{}{
			"code":    http.StatusBadRequest,
			"status":  "400 Bed Request",
			"message": fmt.Sprintf("Unknown messageType:%s", messageType),
			"request": json.Get(data),
		}
		return
	}

}

func (controller *ApiController) addThing(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	controller.subscriptionThings[id] = thing
	for propName, _ := range thing.Properties {
		m := make(map[string]interface{})
		m["id"] = id
		m["messageType"] = models.ThingModified
		value, err := plugin.GetPropertyValue(id, propName)

		if err != nil {
			m["messageType"] = models.ERROR
			m["data"] = map[string]string{"message": err.Error()}

		} else {
			m["messageType"] = models.PropertyStatus
			m["data"] = map[string]interface{}{propName: value}

		}
		controller.sendMessage(m)
	}
}

func (controller *ApiController) onConnected(device *addon.Device, connected bool) {

	t := controller.subscriptionThings[device.ID]
	if t == nil {
		return
	}
	data := make(map[string]interface{})
	data["id"] = t.ID
	data["messageType"] = models.ThingModified
	data["data"] = connected
	controller.sendMessage(data)
}

func (controller *ApiController) onModified(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	t := controller.subscriptionThings[id]
	if t == nil {
		return
	}
	data := make(map[string]interface{})
	data["id"] = t.ID
	data["messageType"] = models.ThingModified
	controller.sendMessage(data)
}

func (controller *ApiController) onRemoved(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	t := controller.subscriptionThings[id]
	if t == nil {
		return
	}
	controller.sendMessage(struct {
		ID          string      `json:"id"`
		MessageType string      `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		ID:          id,
		MessageType: models.ThingRemoved,
		Data: struct {
		}{},
	})
}

func (controller *ApiController) onPropertyChanged(data []byte) {

	deviceId := json.Get(data, "deviceId").ToString()
	name := json.Get(data, "name").ToString()
	v := json.Get(data, "value").GetInterface()
	t := controller.subscriptionThings[deviceId]
	if t == nil {
		return
	}
	var m = map[string]interface{}{name: v}

	controller.sendMessage(struct {
		ID          string      `json:"id"`
		MessageType string      `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		ID:          deviceId,
		MessageType: models.PropertyStatus,
		Data:        m,
	})
}

func (controller *ApiController) sendData(message interface{}) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	data, _ := json.MarshalIndent(&message, "", " ")
	log.Info("things container websocket send message: %s \t\n", string(data))
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, data)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}

func (controller *ApiController) onError(err error) {
	log.Info("websocket err: %s", err.Error())
	controller.done <- struct{}{}
}

func (controller *ApiController) close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	_ = bus.Unsubscribe(util.PropertyChanged, controller.onPropertyChanged)
	_ = bus.Unsubscribe(util.CONNECTED, controller.onConnected)
	_ = bus.Unsubscribe(util.MODIFIED, controller.onModified)
	_ = bus.Unsubscribe(util.ThingRemoved, controller.onRemoved)
	_ = controller.ws.Close()
	return
}

func (controller *ApiController) sendMessage(data map[string]interface{}) {
	d, _ := json.MarshalIndent(&data, "", " ")
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, d)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}
