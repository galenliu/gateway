package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	AddonManager "github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"net/http"
)

type Map = map[string]interface{}

//type WsHandler struct {
//	thingId              string
//	Container            *models.Things
//	ws                   *websocket.Conn
//	locker               *sync.Mutex
//	done                 chan struct{}
//	subscriptionThings   map[string]*models.Thing
//	subscribedEventNames map[string]bool
//}
//
//func NewWsHandler() *WsHandler {
//	controller := &WsHandler{}
//	controller.subscriptionThings = make(map[string]*models.Thing, 10)
//	controller.subscribedEventNames = make(map[string]bool)
//	controller.locker = new(sync.Mutex)
//	controller.Container = models.NewThings()
//	controller.done = make(chan struct{})
//	return controller
//}

func handleWebsocket(c *websocket.Conn) {
	log.Info("handler websocket....")
	thingId, _ := c.Locals("thingId").(string)
	websocketHandler(c, thingId)
}

func websocketHandler(c *websocket.Conn, thingId string) {
	Things := models.NewThings()
	Actions := models.NewActions()
	subscribedEventNames := make(map[string]bool)
	thingCleanups := map[string]func(){}
	var done = make(chan struct{})

	sendMessage := func(data map[string]interface{}) {
		d, _ := json.MarshalIndent(&data, "", " ")
		writeErr := c.WriteMessage(websocket.TextMessage, d)
		if writeErr != nil {
			return
		}
	}

	onEvent := func(event *models.Event) {
		if thingId != "" && event.ThingId != thingId {
			return
		}
		_, ok := subscribedEventNames[event.Name]
		if !ok {
			return
		}
		m := make(map[string]interface{})
		m["id"] = event.ThingId
		m["messageType"] = util.EVENT
		m["data"] = map[string]interface{}{
			event.Name: event.GetDescription(),
		}
		sendMessage(m)
	}

	onActionStatus := func(action *models.Action) {
		if action.ThingId != "" && action.ThingId != thingId {
			return
		}
		m := make(map[string]interface{})
		m["messageType"] = util.ActionStatus
		if action.ThingId != "" {
			m["id"] = action.ThingId
		}
		m["data"] = map[string]interface{}{
			action.Name: action.GetDescription(),
		}
		sendMessage(m)
	}

	addThing := func(thing *models.Thing) {

		onConnected := func(connected bool) {
			m := make(map[string]interface{})
			m["id"] = thing.GetID()
			m["messageType"] = util.CONNECTED
			m["data"] = connected
			sendMessage(m)
		}

		onRemoved := func() {
			f, ok := thingCleanups[thing.GetID()]
			if ok {
				f()
				delete(thingCleanups, thing.GetID())
			}
			m := make(map[string]interface{})
			m["messageType"] = util.ThingRemoved
			m["id"] = thing.GetID()
			m["data"] = struct{}{}
			sendMessage(m)
		}

		onModified := func() {
			m := make(map[string]interface{})
			m["id"] = thing.GetID()
			m["messageType"] = util.ThingModified
			m["data"] = struct{}{}
			sendMessage(m)

		}

		thing.Subscribe(util.CONNECTED, onConnected)
		thing.Subscribe(util.ThingRemoved, onRemoved)
		thing.Subscribe(util.ThingModified, onModified)
		thing.Subscribe(util.EVENT, onEvent)

		thingCleanup := func() {
			thing.Unsubscribe(util.CONNECTED, onConnected)
			thing.Unsubscribe(util.ThingRemoved, onRemoved)
			thing.Unsubscribe(util.ThingModified, onModified)
			thing.Unsubscribe(util.EVENT, onEvent)
		}

		thingCleanups[thing.GetID()] = thingCleanup

		m := make(map[string]interface{})
		m["id"] = thing.GetID()
		m["messageType"] = util.PropertyStatus
		propertyValues := make(map[string]interface{})
		for propName, prop := range thing.Properties {
			value, err := AddonManager.GetPropertyValue(thing.GetID(), propName)
			prop.SetCachedValue(value)
			if err != nil {
				continue
			} else {
				propertyValues[propName] = prop.Value
			}
		}
		m["data"] = propertyValues
		sendMessage(m)
	}

	onThingAdded := func(thing *models.Thing) {
		m := make(map[string]interface{})
		m["id"] = thing.GetID()
		m["messageType"] = util.ThingAdded
		m["data"] = struct{}{}
		sendMessage(m)
		addThing(thing)
	}

	onPropertyChanged := func(data []byte) {
		deviceId := json.Get(data, "deviceId").ToString()
		if thingId != "" && thingId != deviceId {
			return
		}
		name := json.Get(data, "name").ToString()
		v := json.Get(data, "value").GetInterface()
		if name == "" || v == nil {
			return
		}
		m := make(Map)
		m["id"] = deviceId
		m["messageType"] = util.PropertyStatus
		m["data"] = map[string]interface{}{
			name: v,
		}
		sendMessage(m)
	}

	if thingId != "" {
		m := make(map[string]interface{})
		t := Things.GetThing(thingId)
		if t == nil {
			m["messageType"] = util.ERROR
			m["data"] = map[string]interface{}{
				"code":    http.StatusBadRequest,
				"status":  "400 Bed Request",
				"message": fmt.Sprintf("Thing(%s) not found", thingId),
			}
			sendMessage(m)
			return
		}
		addThing(t)
	} else {
		for _, t := range Things.GetThings() {
			addThing(t)
		}
	}

	AddonManager.Subscribe(util.PropertyChanged, onPropertyChanged)
	Things.Subscribe(util.ThingAdded, onThingAdded)
	Actions.Subscribe(util.ActionStatus, onActionStatus)

	clearFunc := func() {
		Things.Unsubscribe(util.ThingAdded, onThingAdded)
		Actions.Unsubscribe(util.ActionStatus, onActionStatus)
		AddonManager.Unsubscribe(util.PropertyChanged, onPropertyChanged)
	}

	defer clearFunc()

	handleMessage := func(bytes []byte) {
		var sendError = func(code int, status string, message string) {
			m := make(map[string]interface{})
			m["messageType"] = util.ERROR
			m["data"] = map[string]interface{}{
				"code":    code,
				"status":  status,
				"message": message,
			}
			sendMessage(m)
		}

		id := json.Get(bytes, "id").ToString()
		if id == "" {
			id = thingId
		}
		device := AddonManager.GetDevice(id)
		messageType := json.Get(bytes, "messageType").ToString()
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

		m := make(map[string]interface{})
		switch messageType {
		case util.SetProperty:
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
				subscribedEventNames[eventName] = true
			}
			return

		case util.RequestAction:
			var actionNames map[string]interface{}
			json.Get(bytes, "data").ToVal(&actionNames)
			for actionName := range actionNames {
				var actionParams map[string]interface{}
				json.Get(bytes, "data", actionName, "input").ToVal(&actionParams)
				th := Things.GetThing(id)
				action := models.NewAction(actionName, actionParams, th)
				Things.Actions.Add(action)
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

	for {
		select {
		case <-done:
			return
		default:
			_, data, readErr := c.ReadMessage()
			if readErr != nil {
				log.Info("websocket disconnected err: %s", readErr.Error())
				done <- struct{}{}
				return
			}
			if data != nil {
				handleMessage(data)
			}

		}

	}

}

//func (controller *WsHandler) handlerConn() {
//
//	if controller.thingId != "" {
//		m := make(map[string]interface{})
//		t := controller.Container.GetThing(controller.thingId)
//		if t == nil {
//			m["messageType"] = util.ERROR
//			m["data"] = map[string]interface{}{
//				"code":    http.StatusBadRequest,
//				"status":  "400 Bed Request",
//				"message": fmt.Sprintf("Thing(%s) not found", controller.thingId),
//			}
//			controller.sendMessage(m)
//			return
//		}
//		controller.addThing(t)
//	} else {
//		for _, t := range controller.Container.GetThings() {
//			controller.addThing(t)
//		}
//	}
//
//	_ = bus.Subscribe(util.PropertyChanged, controller.onPropertyChanged)
//	_ = bus.Subscribe(util.CONNECTED, controller.onConnected)
//	_ = bus.Subscribe(util.MODIFIED, controller.onModified)
//	_ = bus.Subscribe(util.ThingRemoved, controller.onRemoved)
//	_ = bus.Subscribe(util.ThingAdded, controller.onThingAdded)
//
//	defer func() {
//		_ = bus.Unsubscribe(util.PropertyChanged, controller.onPropertyChanged)
//		_ = bus.Unsubscribe(util.CONNECTED, controller.onConnected)
//		_ = bus.Unsubscribe(util.MODIFIED, controller.onModified)
//		_ = bus.Unsubscribe(util.ThingRemoved, controller.onRemoved)
//		_ = bus.Unsubscribe(util.ThingAdded, controller.onThingAdded)
//		if controller.ws != nil {
//			_ = controller.ws.Close()
//		}
//	}()
//
//	for {
//		select {
//		case <-controller.done:
//			return
//		default:
//			_, data, readErr := controller.ws.ReadMessage()
//			if readErr != nil {
//				log.Info("websocket disconnected err: %s", readErr.Error())
//				controller.done <- struct{}{}
//				return
//			}
//			if data != nil {
//				controller.handleMessage(data)
//			}
//
//		}
//
//	}
//
//}
//
//func (controller *WsHandler) handleMessage(bytes []byte) {
//
//	var sendError = func(code int, status string, message string) {
//		m := make(map[string]interface{})
//		m["messageType"] = util.ERROR
//		m["data"] = map[string]interface{}{
//			"code":    code,
//			"status":  status,
//			"message": message,
//		}
//		controller.sendMessage(m)
//	}
//
//	id := json.Get(bytes, "id").ToString()
//	if id == "" {
//		id = controller.thingId
//	}
//	device := AddonManager.GetDevice(id)
//	messageType := json.Get(bytes, "messageType").ToString()
//	m := make(map[string]interface{})
//
//	if id == "" {
//		sendError(400, "400 Bed Request", "Missing thing id")
//		return
//	}
//	if device == nil {
//		sendError(400, "400 Bed Request", "device can not found")
//		return
//	}
//	if messageType == "" {
//		sendError(400, "400 Bed Request", "messageType err")
//		return
//	}
//
//	switch messageType {
//	case util.SetProperty:
//		var propertyMap map[string]interface{}
//		json.Get(bytes, "data").ToVal(&propertyMap)
//		for propName, value := range propertyMap {
//			_, setErr := AddonManager.SetProperty(device.GetID(), propName, value)
//			if setErr != nil {
//				m["messageType"] = util.ERROR
//				m["bytes"] = map[string]interface{}{
//					"code":    http.StatusBadRequest,
//					"status":  "400 Bed Request",
//					"message": setErr.Error(),
//					"request": bytes,
//				}
//			}
//		}
//		return
//	case util.AddEventSubscription:
//		var eventsName []string
//		json.Get(bytes, "data").ToVal(&eventsName)
//		for _, eventName := range eventsName {
//			controller.subscribedEventNames[eventName] = true
//		}
//		return
//
//	case util.RequestAction:
//		var actionNames map[string]interface{}
//		json.Get(bytes, "data").ToVal(&actionNames)
//		for actionName, _ := range actionNames {
//			var actionParams map[string]interface{}
//			json.Get(bytes, "data", actionName, "input").ToVal(&actionParams)
//			th := controller.Container.GetThing(id)
//			action := models.NewAction(actionName, actionParams, th)
//			controller.Container.Actions.Add(action)
//			err := AddonManager.RequestAction(id, action.ID, actionName, actionParams)
//			if err != nil {
//				sendError(400, "400 Bad Request", err.Error())
//			}
//		}
//	default:
//		sendError(400, "400 Bed Request", fmt.Sprintf("Unknown messageType:%s", messageType))
//		return
//	}
//}
//
//func (controller *WsHandler) addThing(thing *models.Thing) {
//
//	sl := strings.Split(thing.ID, "/")
//	id := sl[len(sl)-1]
//	controller.subscriptionThings[id] = thing
//
//	m := make(map[string]interface{})
//	m["id"] = id
//	for propName, _ := range thing.Properties {
//		value, err := AddonManager.GetPropertyValue(id, propName)
//		if err != nil {
//			m["messageType"] = util.PropertyStatus
//			m["data"] = map[string]interface{}{"message": "property set err"}
//			continue
//		} else {
//			m["data"] = map[string]interface{}{propName: value}
//		}
//		controller.sendMessage(m)
//	}
//
//	controller.onConnected(id, thing.IsConnected())
//}
//
//func (controller *WsHandler) onThingAdded(thing *models.Thing) {
//
//	sl := strings.Split(thing.ID, "/")
//	id := sl[len(sl)-1]
//	m := make(map[string]interface{})
//	m["id"] = id
//	m["messageType"] = util.ThingAdded
//	m["data"] = struct{}{}
//	controller.sendMessage(m)
//	controller.addThing(thing)
//}
//
//func (controller *WsHandler) onConnected(deviceId string, connected bool) {
//
//	_, ok := controller.subscriptionThings[deviceId]
//	if !ok {
//		return
//	}
//	m := make(map[string]interface{})
//	m["id"] = deviceId
//	m["messageType"] = util.CONNECTED
//	m["data"] = connected
//	controller.sendMessage(m)
//}
//
//func (controller *WsHandler) onModified(thing *models.Thing) {
//
//	sl := strings.Split(thing.ID, "/")
//	id := sl[len(sl)-1]
//	t := controller.subscriptionThings[id]
//	if t == nil {
//		return
//	}
//	m := make(map[string]interface{})
//	m["id"] = t.ID
//	m["messageType"] = util.ThingModified
//	controller.sendMessage(m)
//}
//
//func (controller *WsHandler) onActionStatus(action *models.Action) {
//	if action.ThingId == "" {
//		return
//	}
//	for _, th := range controller.subscriptionThings {
//		if th.ID == action.ThingId {
//			m := make(map[string]interface{})
//			m["id"] = action.ThingId
//			m["messageType"] = util.ActionStatus
//			m["data"] = map[string]interface{}{
//				action.Name: action.MarshalJson(),
//			}
//			controller.sendMessage(m)
//
//		}
//	}
//}
//
//func (controller *WsHandler) onEvent(event *models.Event) {
//	if !controller.subscribedEventNames[event.Name] {
//		return
//	}
//	m := make(map[string]interface{})
//	m["id"] = event.GetThingId()
//	m["messageType"] = util.EVENT
//	m["data"] = map[string]interface{}{event.Name: event.MarshalJson()}
//
//}
//
//func (controller *WsHandler) onRemoved(thing *models.Thing) {
//
//	sl := strings.Split(thing.ID, "/")
//	id := sl[len(sl)-1]
//	t := controller.subscriptionThings[id]
//	if t == nil {
//		return
//	}
//	m := make(map[string]interface{})
//	m["id"] = id
//	m["messageType"] = util.ThingRemoved
//	m["data"] = map[string]interface{}{}
//	delete(controller.subscriptionThings, id)
//	controller.sendMessage(m)
//}
//
//func (controller *WsHandler) onPropertyChanged(data []byte) {
//
//	deviceId := json.Get(data, "deviceId").ToString()
//	name := gjson.GetBytes(data, "name").String()
//	v := gjson.GetBytes(data, "value").Value()
//	t := controller.subscriptionThings[deviceId]
//	if t == nil {
//		return
//	}
//	m := make(map[string]interface{})
//	m["id"] = deviceId
//	m["messageType"] = util.PropertyStatus
//	m["data"] = map[string]interface{}{
//		name: v,
//	}
//	controller.sendMessage(m)
//}
//
//func (controller *WsHandler) sendData(message interface{}) {
//	controller.locker.Lock()
//	defer controller.locker.Unlock()
//	data, _ := json.MarshalIndent(&message, "", " ")
//	log.Info("things container websocket send message: %s \t\n", string(data))
//	writeErr := controller.ws.WriteMessage(websocket.TextMessage, data)
//	if writeErr != nil {
//		controller.onError(writeErr)
//	}
//}
//
//func (controller *WsHandler) onError(err error) {
//	log.Info("websocket err: %s", err.Error())
//	controller.done <- struct{}{}
//}
//
//func (controller *WsHandler) sendMessage(data map[string]interface{}) {
//	if controller.ws == nil {
//		log.Info("websocket nil")
//		return
//	}
//	d, _ := json.MarshalIndent(&data, "", " ")
//	writeErr := controller.ws.WriteMessage(websocket.TextMessage, d)
//	if writeErr != nil {
//		controller.onError(writeErr)
//	}
//}
