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

type ThingsWebsocketHandler struct {
	thingId            string
	Container          *models.Things
	ws                 *websocket.Conn
	locker             *sync.Mutex
	done               chan struct{}
	subscriptionThings map[string]*thing.Thing
}

func NewThingsWebsocketController() *ThingsWebsocketHandler {
	controller := &ThingsWebsocketHandler{}
	controller.subscriptionThings = make(map[string]*thing.Thing, 10)
	controller.locker = new(sync.Mutex)
	controller.Container = models.NewThings()
	controller.done = make(chan struct{})
	return controller
}

func handleWebsocket(c *websocket.Conn) {

	if !c.Locals("websocket").(bool) {
		return
	}
	controller := NewThingsWebsocketController()
	controller.ws = c
	controller.thingId = c.Params("thingId")
	if controller.thingId != "" {
		t := controller.Container.GetThing(controller.thingId)
		if t == nil {
			log.Info(fmt.Sprintf("Thing(%s) not found", controller.thingId))
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
					Message: fmt.Sprintf("Thing(%s) not found", controller.thingId),
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

func (controller *ThingsWebsocketHandler) handleMessage(data []byte) {

	id := json.Get(data, "id").ToString()
	if id == "" {
		id = controller.thingId
	}
	// get devices form addon
	device := plugin.GetDevice(id)
	messageType := json.Get(data, "messageType").ToString()

	if id == "" || device == nil || messageType == "" {
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

	switch messageType {
	case models.SetProperty:
		var propertyMap map[string]interface{}
		json.Get(data, "data").ToVal(&propertyMap)

		for propName, value := range propertyMap {
			p, setErr := plugin.SetProperty(device.ID, propName, value)
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
							Value: p.Value,
						},
						Request: json.Get(data),
					},
				})

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

func (controller *ThingsWebsocketHandler) addThing(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	controller.subscriptionThings[id] = thing
	for propName, _ := range thing.Properties {
		value, err := plugin.GetPropertyValue(id, propName)
		if err != nil {
			controller.sendMessage(struct {
				ID          string      `json:"id"`
				MessageType string      `json:"messageType"`
				Data        interface{} `json:"data"`
			}{
				ID:          id,
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
				ID:          id,
				MessageType: models.PropertyStatus,
				Data:        map[string]interface{}{propName: value},
			})
		}
	}
}

func (controller *ThingsWebsocketHandler) onConnected(device *addon.Device, connected bool) {

	t := controller.subscriptionThings[device.ID]
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

func (controller *ThingsWebsocketHandler) onModified(thing *thing.Thing) {

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
		MessageType: models.ThingModified,
		Data: struct {
		}{},
	})
}

func (controller *ThingsWebsocketHandler) onRemoved(thing *thing.Thing) {

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

func (controller *ThingsWebsocketHandler) onPropertyChanged(property *addon.Property) {

	t := controller.subscriptionThings[property.DeviceId]
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

func (controller *ThingsWebsocketHandler) sendMessage(message interface{}) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	data, _ := json.MarshalIndent(&message, "", " ")
	log.Info("things container websocket send message: %s \t\n", string(data))
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, data)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}

func (controller *ThingsWebsocketHandler) onError(err error) {
	log.Info("websocket err: %s", err.Error())
	controller.done <- struct{}{}
}

func (controller *ThingsWebsocketHandler) close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	_ = bus.Unsubscribe(util.PropertyChanged, controller.onPropertyChanged)
	_ = bus.Unsubscribe(util.CONNECTED, controller.onConnected)
	_ = bus.Unsubscribe(util.MODIFIED, controller.onModified)
	_ = bus.Unsubscribe(util.ThingRemoved, controller.onRemoved)
	_ = controller.ws.Close()
	return
}
