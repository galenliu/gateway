package controllers

import (
	"addon"
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	AddonManager "gateway/plugin"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"strings"

	"sync"
)

type WsHandler struct {
	thingId              string
	Container            *models.Things
	ws                   *websocket.Conn
	locker               *sync.Mutex
	done                 chan struct{}
	subscriptionThings   map[string]*thing.Thing
	subscribedEventNames map[string]bool
}

func NewWsHandler() *WsHandler {
	controller := &WsHandler{}
	controller.subscriptionThings = make(map[string]*thing.Thing, 10)
	controller.subscribedEventNames = make(map[string]bool)
	controller.locker = new(sync.Mutex)
	controller.Container = models.NewThings()
	controller.done = make(chan struct{})
	return controller
}

func handleWebsocket(c *websocket.Conn) {
	if !c.Locals("websocket").(bool) {
		return
	}
	log.Info("handler websocket....")
	controller := NewWsHandler()
	controller.thingId, _ = c.Locals("thingId").(string)
	controller.ws = c
	controller.handlerConn()

}

func (controller *WsHandler) handlerConn() {

	if controller.thingId != "" {
		m := make(map[string]interface{})
		t := controller.Container.GetThing(controller.thingId)
		if t == nil {
			m["messageType"] = util.ERROR
			m["data"] = map[string]interface{}{
				"code":    http.StatusBadRequest,
				"status":  "400 Bed Request",
				"message": fmt.Sprintf("Thing(%s) not found", controller.thingId),
			}
			controller.sendMessage(m)
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
				controller.handleMessage(data)
			}

		}

	}

}

func (controller *WsHandler) handleMessage(bytes []byte) {

	var sendError = func(code int, status string, message string) {
		m := make(map[string]interface{})
		m["messageType"] = util.ERROR
		m["data"] = map[string]interface{}{
			"code":    code,
			"status":  status,
			"message": message,
		}
		controller.sendMessage(m)
	}

	id := json.Get(bytes, "id").ToString()
	if id == "" {
		id = controller.thingId
	}
	device := AddonManager.GetDevice(id)
	messageType := json.Get(bytes, "messageType").ToString()
	m := make(map[string]interface{})

	if id == "" {
		sendError(400, "400 Bed Request", "Missing thing id")
		return
	}
	if device == nil {
		sendError(400, "400 Bed Request", "device can not found")
		return
	}
	if messageType == "" {
		sendError(400, "400 Bed Request", "messageType err")
		return
	}

	switch messageType {
	case models.SetProperty:
		var propertyMap map[string]interface{}
		json.Get(bytes, "data").ToVal(&propertyMap)
		for propName, value := range propertyMap {
			_, setErr := AddonManager.SetProperty(device.GetID(), propName, value)
			if setErr != nil {
				m["messageType"] = util.ERROR
				m["bytes"] = map[string]interface{}{
					"code":    http.StatusBadRequest,
					"status":  "400 Bed Request",
					"message": setErr.Error(),
					"request": bytes,
				}
			}
		}
		return
	case util.AddEventSubscription:
		var eventsName []string
		json.Get(bytes, "data").ToVal(&eventsName)
		for _, eventName := range eventsName {
			controller.subscribedEventNames[eventName] = true
		}
		return

	case util.RequestAction:
		var actionNames map[string]interface{}
		json.Get(bytes, "data").ToVal(&actionNames)
		for actionName, _ := range actionNames {
			var actionParams map[string]interface{}
			json.Get(bytes, "data", actionName, "input").ToVal(&actionParams)
			th := controller.Container.GetThing(id)
			action := thing.NewAction(actionName, actionParams, th)
			controller.Container.Actions.Add(action)
			err := AddonManager.RequestAction(id, action.ID, actionName, actionParams)
			if err != nil {
				sendError(400, "400 Bad Request", err.Error())
			}
		}

	default:
		sendError(400, "400 Bed Request", fmt.Sprintf("Unknown messageType:%s", messageType))
		return
	}

}

func (controller *WsHandler) addThing(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	controller.subscriptionThings[id] = thing
	for propName, _ := range thing.Properties {
		m := make(map[string]interface{})
		m["id"] = id
		m["messageType"] = models.ThingModified
		value, err := AddonManager.GetPropertyValue(id, propName)

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

func (controller *WsHandler) onConnected(device *addon.Device, connected bool) {

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

func (controller *WsHandler) onModified(thing *thing.Thing) {

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

func (controller *WsHandler) onActionStatus(action *thing.Action) {
	if action.ThingId == "" {
		return
	}
	for _, th := range controller.subscriptionThings {
		if th.ID == action.ThingId {
			m := make(map[string]interface{})
			m["id"] = action.ThingId
			m["messageType"] = util.ActionStatus
			m["data"] = map[string]interface{}{
				action.Name: action.GetDescription(),
			}
			controller.sendMessage(m)

		}
	}
}

func (controller *WsHandler) onEvent(event *thing.Event) {
	if !controller.subscribedEventNames[event.Name] {
		return
	}
	m := make(map[string]interface{})
	m["id"] = event.GetThingId()
	m["messageType"] = util.EVENT
	m["data"] = map[string]interface{}{event.Name: event.GetDescription()}

}

func (controller *WsHandler) onRemoved(thing *thing.Thing) {

	sl := strings.Split(thing.ID, "/")
	id := sl[len(sl)-1]
	t := controller.subscriptionThings[id]
	if t == nil {
		return
	}
	m := make(map[string]interface{})
	m["id"] = id
	m["messageType"] = models.ThingRemoved
	m["data"] = map[string]interface{}{}
	controller.sendMessage(m)
}

func (controller *WsHandler) onPropertyChanged(data []byte) {

	deviceId := json.Get(data, "deviceId").ToString()
	name := json.Get(data, "name").ToString()
	v := json.Get(data, "value").GetInterface()
	t := controller.subscriptionThings[deviceId]
	if t == nil {
		return
	}

	m := make(map[string]interface{})
	m["id"] = deviceId
	m["messageType"] = models.PropertyStatus
	m["data"] = map[string]interface{}{
		name: v,
	}
	controller.sendMessage(m)
}

func (controller *WsHandler) sendData(message interface{}) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	data, _ := json.MarshalIndent(&message, "", " ")
	log.Info("things container websocket send message: %s \t\n", string(data))
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, data)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}

func (controller *WsHandler) onError(err error) {
	log.Info("websocket err: %s", err.Error())
	controller.done <- struct{}{}
}

func (controller *WsHandler) close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	_ = bus.Unsubscribe(util.PropertyChanged, controller.onPropertyChanged)
	_ = bus.Unsubscribe(util.CONNECTED, controller.onConnected)
	_ = bus.Unsubscribe(util.MODIFIED, controller.onModified)
	_ = bus.Unsubscribe(util.ThingRemoved, controller.onRemoved)
	_ = controller.ws.Close()
	return
}

func (controller *WsHandler) sendMessage(data map[string]interface{}) {
	if controller.ws == nil {
		log.Info("websocket nil")
		return
	}
	d, _ := json.MarshalIndent(&data, "", " ")
	writeErr := controller.ws.WriteMessage(websocket.TextMessage, d)
	if writeErr != nil {
		controller.onError(writeErr)
	}
}
