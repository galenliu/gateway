package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

type wsContainer struct {
	thingId   string
	container things.Container
}

func handleWebsocket(container things.Container, log logging.Logger) func(ws *websocket.Conn) {
	handler := func(ws *websocket.Conn) {
		defer ws.Close()
		//subscribedEventNames := make(map[string]bool)
		thingCleanups := make(map[string]func())
		onClose := func() {
			if thingCleanups != nil && len(thingCleanups) > 0 {
				for _, f := range thingCleanups {
					f()
				}
			}
		}
		defer onClose()
		thingId := ws.Params("thingId")

		log.Infof("websocket connection")

		write := func(id string, messageType string, data any) {
			if ws != nil {
				err := ws.WriteJSON(struct {
					Id          string `json:"id"`
					MessageType string `json:"messageType"`
					Data        any    `json:"data,omitempty"`
				}{
					id,
					messageType,
					data,
				})
				if err != nil {
					log.Error("websocket error: messageType:%s  err : %s", messageType, err.Error())
					return
				}
			}
		}

		addThing := func(t *things.Thing) {
			onThingConnected := func(message topic.ThingConnectedMessage) {
				if message.ThingId != t.GetId() {
					return
				}
				write(t.GetId(), constant.Connected, message.Connected)
			}
			onThingConnected(topic.ThingConnectedMessage{
				ThingId:   t.GetId(),
				Connected: t.Connected,
			})
			removeConnectedFunc := container.Subscribe(topic.ThingConnected, onThingConnected)

			onThingRemoved := func(message topic.ThingRemovedMessage) {
				if message.ThingId != t.GetId() {
					return
				}
				//f, ok := thingCleanups[t.GetId()]
				//if ok {
				//	f()
				//}
				delete(thingCleanups, t.GetId())
				write(t.GetId(), constant.ThingRemoved, nil)
			}
			removeRemovedFunc := container.Subscribe(topic.ThingRemoved, onThingRemoved)

			onThingModified := func(message topic.ThingModifyMessage) {
				if message.ThingId != t.GetId() {
					return
				}
				write(t.GetId(), constant.ThingModified, nil)
			}
			removeModifiedFunc := container.Subscribe(topic.ThingModify, onThingModified)

			onEvent := func(e *events.Event) {
				write(t.GetId(), constant.Event, struct {
					Name  string        `json:"name"`
					Event *events.Event `json:"events"`
				}{Name: e.GetName(), Event: e})
			}
			removeThingEventStatusFunc := container.Subscribe(topic.ThingEvent, onEvent)

			onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
				if message.ThingId != t.GetId() {
					return
				}
				write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
			}
			removePropertyChangedFunc := container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)

			onActionStatus := func(thingId string, action *actions.ActionDescription) {
				if t.GetId() != thingId {
					return
				}
				write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
			}
			removeActionStatusFunc := container.Subscribe(topic.ThingActionStatus, onActionStatus)

			clean := func() {
				removeConnectedFunc()
				removeRemovedFunc()
				removeModifiedFunc()
				removePropertyChangedFunc()
				removeActionStatusFunc()
				removeThingEventStatusFunc()
			}
			thingCleanups[t.GetId()] = clean

			var data = make(map[string]any)
			for n, _ := range t.Properties {
				v, err := t.GetPropertyValue(n)
				if err == nil {
					data[n] = v
				}
				write(t.GetId(), constant.PropertyStatus, data)
			}
		}

		handleMessage := func(data []byte) {
			log.Infof("received message: %v", data)
		}

		onThingAdded := func(message topic.ThingAddedMessage) {
			if ws == nil {
				return
			}
			err := ws.WriteJSON(map[string]any{
				"id":          message.ThingId,
				"messageType": constant.ThingAdded,
				"data":        struct{}{},
			})
			if err != nil {
				log.Error("websocket send %s message err : %s", constant.ThingAdded, err.Error())
			}
			thing := container.GetThing(message.ThingId)
			if thing != nil {
				addThing(thing)
			}
		}

		if thingId == "" {
			for _, thing := range container.GetThings() {
				addThing(thing)
			}
			container.Subscribe(topic.ThingAdded, onThingAdded)
		} else {
			t := container.GetThing(thingId)
			if t == nil {
				return
			}
			addThing(t)
		}

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if message == nil {
				log.Infof("received message: %s", message)
				handleMessage(message)
			}
		}
	}
	return handler
}
