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
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"net/http"
	"sync"
)

type ThingsWebsocketController struct {
	thingId            string
	Container          *models.Things
	ws                 *websocket.Conn
	locker             *sync.Mutex
	subscriptionThings []*thing.Thing
}

func NewThingsWebsocketController() *ThingsWebsocketController {
	handler := &ThingsWebsocketController{}

	return handler
}

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebsocket(c *gin.Context, things *models.Things) {

	//  websocket upgrade
	log.Info("handle websocket connection:", c.Request)
	conn, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	controller := NewThingsWebsocketController()
	controller.ws = conn
	controller.thingId = c.Param("thingId")
	controller.locker = new(sync.Mutex)
	controller.Container = things
	bus.Subscribe(util.PropertyChanged, controller.onPropertyChanged)
	bus.Subscribe(util.CONNECTED, controller.onConnected)
	bus.Subscribe(util.MODIFIED, controller.onModified)
	bus.Subscribe(util.ThingRemoved, controller.onRemoved)

	if controller.thingId != "" {
		t := controller.Container.GetThing(controller.thingId)
		if t == nil {
			log.Info(fmt.Sprintf("THing(%s) not found", controller.thingId))
			controller.sendMessage(struct {
				messageType string
				data        interface{}
			}{
				messageType: thing.Error,
				data: struct {
					Code    int
					Status  string
					Message string `json:"message"`
				}{
					Code:    400,
					Status:  "404 NOT FOUND",
					Message: fmt.Sprintf("THing(%s) not found", controller.thingId),
				},
			})
			controller.close()
			return
		}
		controller.addThing(t)
	} else {
		for _, t := range controller.Container.GetThings() {
			controller.addThing(t)
		}
	}
	go controller.handleWebsocket()

}

func (controller *ThingsWebsocketController) handleWebsocket() {
	for {
		_, data, readErr := controller.ws.ReadMessage()
		if data == nil {
			break
		}
		if readErr != nil {
			log.Info("websocket disconnected", readErr.Error())
			controller.close()
			return
		}
		go controller.handleMessage(data)
	}
}

func (controller *ThingsWebsocketController) handleMessage(data []byte) {

	id := json.Get(data, "id").ToString()
	if id == "" {
		id = controller.thingId
	}

	// get thing id form request message
	if id == "" {
		controller.sendMessage(struct {
			MessageType string `json:"messageType"`
			Data        interface{}
		}{
			MessageType: thing.Error,
			Data: struct {
				Code    int    `json:"code"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Code:    http.StatusBadRequest,
				Status:  "400 Bed Request",
				Message: "Messing thing id",
			},
		})

		return
	}

	// get devices form addon
	device, err := plugin.FindDevice(id)

	if err != nil {
		controller.sendMessage(struct {
			MessageType string      `json:"messageType"`
			Data        interface{} `json:"data"`
		}{
			MessageType: thing.Error,
			Data: struct {
				Code    int         `json:"code"`
				Status  string      `json:"status"`
				Message string      `json:"message"`
				Request interface{} `json:"request"`
			}{
				Code:    http.StatusBadRequest,
				Status:  "400 Bed Request",
				Message: fmt.Sprintf("thing id(%s) not found", id),
				Request: json.Get(data),
			},
		})

		return
	}

	messageType := json.Get(data, "messageType").ToString()

	switch messageType {
	case models.SetProperty:
		var propertyMap map[string]interface{}
		json.Get(data, "data").ToVal(&propertyMap)
		for propName, value := range propertyMap {
			setErr := plugin.SetProperValue(device.ID, propName, value)
			if setErr != nil {
				controller.sendMessage(struct {
					MessageType string      `json:"messageType"`
					Data        interface{} `json:"data"`
				}{
					MessageType: thing.Error,
					Data: struct {
						Code    int         `json:"code"`
						Status  string      `json:"status"`
						Message string      `json:"message"`
						Request interface{} `json:"request"`
					}{
						Code:    http.StatusBadRequest,
						Status:  "400 Bed Request",
						Message: setErr.Error(),
						Request: json.Get(data),
					},
				})
			} else {

			}
		}
		return

	default:
		controller.sendMessage(struct {
			MessageType string `json:"messageType"`
			Data        interface{}
		}{
			MessageType: thing.Error,
			Data: struct {
				Code    int         `json:"code"`
				Status  string      `json:"status"`
				Message string      `json:"message"`
				Request interface{} `json:"request"`
			}{
				Code:    http.StatusBadRequest,
				Status:  "400 Bed Request",
				Message: fmt.Sprintf("Unknown messageType:%s", messageType),
				Request: json.Get(data),
			},
		})
		return
	}

}

func (controller *ThingsWebsocketController) addThing(thing *thing.Thing) {

	controller.subscriptionThings = append(controller.subscriptionThings, thing)

	for propName, _ := range thing.Properties {
		value, err := plugin.GetPropertyValue(thing.ID, propName)
		if err != nil {
			controller.sendMessage(struct {
				ID          string      `json:"id"`
				MessageType string      `json:"messageType"`
				Data        interface{} `json:"data"`
			}{
				ID:          thing.ID,
				MessageType: models.ERROR,
				Data: struct {
					Message string `json:"message"`
				}{Message: err.Error()},
			})
		} else {
			controller.sendMessage(struct {
				ID          string      `json:"id"`
				MessageType string      `json:"messageType"`
				Data        interface{} `json:"data"`
			}{
				ID:          thing.ID,
				MessageType: models.PropertyStatus,
				Data:        map[string]interface{}{propName: value},
			})
		}
	}
}

func (controller *ThingsWebsocketController) onConnected(device *addon.Device, connected bool) {

	t := controller.getThing(device.ID)
	if t == nil {
		return
	}
	controller.sendMessage(struct {
		ID          string `json:"id"`
		MessageType string `json:"messageType"`
		Data        bool   `json:"data"`
	}{
		ID:          t.ID,
		MessageType: models.CONNECTED,
		Data:        connected,
	})

}

func (controller *ThingsWebsocketController) onModified(thing *thing.Thing) {

	t := controller.getThing(thing.ID)
	if t == nil {
		return
	}
	controller.sendMessage(struct {
		ID          string      `json:"id"`
		MessageType string      `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		ID:          t.ID,
		MessageType: models.ThingModified,
		Data: struct {
		}{},
	})
}

func (controller *ThingsWebsocketController) onRemoved(thing *thing.Thing) {

	t := controller.getThing(thing.ID)
	if t == nil {
		return
	}
	controller.sendMessage(struct {
		ID          string      `json:"id"`
		MessageType string      `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		ID:          thing.ID,
		MessageType: models.ThingRemoved,
		Data: struct {
		}{},
	})
}

func (controller *ThingsWebsocketController) onPropertyChanged(property *addon.Property) {

	t := controller.getThing(property.DeviceId)
	if t == nil {
		return
	}
	var data = map[string]interface{}{property.Name: property.Value}
	controller.sendMessage(struct {
		ID          string      `json:"id"`
		MessageType string      `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		ID:          property.DeviceId,
		MessageType: models.PropertyStatus,
		Data:        data,
	})
}

func (controller *ThingsWebsocketController) sendMessage(message interface{}) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	data, _ := json.MarshalIndent(&message, "", " ")
	log.Info("things container websocket send message: %s \t\n", string(data))
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, data)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}

func (controller *ThingsWebsocketController) onError(err error) {
	log.Info("websocket err:", err.Error())
	controller.close()
}

func (controller *ThingsWebsocketController) close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	bus.Unsubscribe(util.PropertyChanged, controller.onPropertyChanged)
	bus.Unsubscribe(util.CONNECTED, controller.onConnected)
	bus.Unsubscribe(util.MODIFIED, controller.onModified)
	bus.Unsubscribe(util.ThingRemoved, controller.onRemoved)
	_ = controller.ws.Close()
}

func (controller *ThingsWebsocketController) getThing(id string) *thing.Thing {
	for _, t := range controller.subscriptionThings {
		if t.ID == id {
			return t
		}
	}
	return nil
}
