package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 2 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func handleWebsocket(container things.Container, log logging.Logger) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		locker := sync.Mutex{}
		isClosed := false
		//subscribedEventNames := make(map[string]bool)
		thingCleanups := make(map[string]func())
		closeWs := func() {
			locker.Lock()
			defer locker.Unlock()
			if !isClosed {
				log.Infof("closed websocket connection")
				ws.Close()
			}
			isClosed = true
			if thingCleanups != nil && len(thingCleanups) > 0 {
				for id, f := range thingCleanups {
					log.Infof("thing %s thingCleanup %v", id, f)
					f()
				}
			}
			log.Infof("websocket close")
		}
		defer closeWs()
		thingId := ws.Params("thingId")

		write := func(id string, messageType string, d ...any) {
			js := struct {
				Id          string `json:"id"`
				MessageType string `json:"messageType"`
				Data        any    `json:"data,omitempty"`
			}{
				id,
				messageType,
				func() any {
					if d == nil {
						return struct{}{}
					}
					return d[0]
				}(),
			}
			log.Infof("write message: %v", js)
			locker.Lock()
			defer locker.Unlock()
			if isClosed {
				return
			}
			if ws.Conn != nil || !isClosed {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				err := ws.WriteJSON(js)
				if err != nil {
					log.Error("websocket error : %s", err.Error())
					return
				}
			}
		}

		//addThing := func(t *things.Thing) {
		//onThingConnected := func(message topic.ThingConnectedMessage) {
		//	if message.ThingId != t.GetId() {
		//		return
		//	}
		//	write(t.GetId(), constant.Connected, message.Connected)
		//}
		//log.Infof("addThing %s", t.Id)
		//write(t.GetId(), constant.Connected, t.Connected)
		//removeConnectedFunc := container.Subscribe(topic.ThingConnected, onThingConnected)

		//onThingRemoved := func(message topic.ThingRemovedMessage) {
		//	if message.ThingId != t.GetId() {
		//		return
		//	}
		//	//f, ok := thingCleanups[t.GetId()]
		//	//if ok {
		//	//	f()
		//	//}
		//	delete(thingCleanups, t.GetId())
		//	write(t.GetId(), constant.ThingRemoved)
		//}
		//removeRemovedFunc := container.Subscribe(topic.ThingRemoved, onThingRemoved)

		//onThingModified := func(message topic.ThingModifyMessage) {
		//	if message.ThingId != t.GetId() {
		//		return
		//	}
		//	write(t.GetId(), constant.ThingModified)
		//}
		//removeModifiedFunc := container.Subscribe(topic.ThingModify, onThingModified)

		//onEvent := func(e *events.Event) {
		//	write(t.GetId(), constant.Event, struct {
		//		Name  string        `json:"name"`
		//		Event *events.Event `json:"events"`
		//	}{Name: e.GetName(), Event: e})
		//}
		//removeThingEventStatusFunc := container.Subscribe(topic.ThingEvent, onEvent)

		//onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
		//	if message.ThingId != t.GetId() {
		//		return
		//	}
		//	write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
		//}
		//removePropertyChangedFunc := container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)

		//onActionStatus := func(thingId string, action *actions.ActionDescription) {
		//	if t.GetId() != thingId {
		//		return
		//	}
		//	write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
		//}
		//removeActionStatusFunc := container.Subscribe(topic.ThingActionStatus, onActionStatus)

		//clean := func() {
		//removeConnectedFunc()
		//removeRemovedFunc()
		//removeModifiedFunc()
		//removePropertyChangedFunc()
		//removeActionStatusFunc()
		//removeThingEventStatusFunc()
		//}
		//thingCleanups[t.GetId()] = clean

		//}

		handleMessage := func(data []byte) {
			log.Infof("received message: %v", data)
		}

		onThingConnected := func(id string) func(message topic.ThingConnectedMessage) {
			return func(message topic.ThingConnectedMessage) {
				if message.ThingId != id {
					return
				}
				write(id, constant.Connected, message.Connected)
			}
		}

		onThingRemoved := func(id string) func(message topic.ThingRemovedMessage) {
			return func(message topic.ThingRemovedMessage) {
				if message.ThingId != id {
					return
				}
				delete(thingCleanups, id)
				write(id, constant.ThingRemoved)

			}
		}

		onThingModified := func(id string) func(message topic.ThingRemovedMessage) {
			return func(message topic.ThingRemovedMessage) {
				if message.ThingId != id {
					return
				}
				write(id, constant.ThingModified)
			}
		}

		onPropertyChanged := func(id string) func(message topic.ThingPropertyChangedMessage) {
			return func(message topic.ThingPropertyChangedMessage) {
				if message.ThingId != id {
					return
				}
				write(id, constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
			}
		}

		onActionStatus := func(id string) func(aid string, action *actions.ActionDescription) {
			return func(aid string, action *actions.ActionDescription) {
				if aid != id {
					return
				}
				write(id, constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
			}
		}

		//onActionStatus := func(thingId string, action *actions.ActionDescription) {
		//	if t.GetId() != thingId {
		//		return
		//	}
		//	write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
		//}
		//removeActionStatusFunc := container.Subscribe(topic.ThingActionStatus, onActionStatus)

		//onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
		//	if message.ThingId != t.GetId() {
		//		return
		//	}
		//	write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
		//}
		//removePropertyChangedFunc := container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)

		onEvent := func(id string) func(e *events.Event) {
			return func(e *events.Event) {
				write(id, constant.Event, struct {
					Name  string        `json:"name"`
					Event *events.Event `json:"events"`
				}{Name: e.GetName(), Event: e})
			}
		}

		add := func(t1 *things.Thing) {
			container.Subscribe(topic.ThingConnected, onThingConnected(t1.GetId()))
			container.Subscribe(topic.ThingRemoved, onThingRemoved(t1.GetId()))
			container.Subscribe(topic.ThingModify, onThingModified(t1.GetId()))
			container.Subscribe(topic.ThingEvent, onEvent(t1.GetId()))
			container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged(t1.GetId()))
			container.Subscribe(topic.ThingActionStatus, onActionStatus(t1.GetId()))
			write(t1.GetId(), constant.Connected, t1.Connected)
			var data = make(map[string]any)
			for n, _ := range t1.Properties {
				v, err := t1.GetPropertyValue(n)
				if err == nil {
					data[n] = v
				}
				write(t1.GetId(), constant.PropertyStatus, data)
			}
		}

		onThingAdded := func(message topic.ThingAddedMessage) {
			locker.Lock()
			locker.Unlock()
			if ws.Conn == nil || isClosed {
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
				add(thing)
			}
		}

		if thingId == "" {
			ts := container.GetThings()
			for _, t1 := range ts {
				add(t1)
			}

			container.Subscribe(topic.ThingAdded, onThingAdded)
		} else {
			t2 := container.GetThing(thingId)
			if t2 == nil {
				return
			}
			add(t2)
		}

		ws.SetReadDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error {
			log.Infof("ipc server pong handler")
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Infof("things websocket read error: %v", err.Error())
				return
			}
			if message == nil {
				log.Infof("received message: %s", message)
				handleMessage(message)
			}
		}
	}
}
