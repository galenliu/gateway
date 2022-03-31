package controllers

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
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

var connectionId int64 = 0

type message struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
	Data        any    `json:"data,omitempty"`
}

type Connection struct {
	*websocket.Conn
	thingId   string
	isClosed  bool
	locker    sync.Mutex
	container things.Container
	log       logging.Logger
	connId    int64
}

func (conn *Connection) write(id string, messageType string, data any) error {
	conn.locker.Lock()
	defer conn.locker.Unlock()
	if conn.Conn != nil && !conn.isClosed {
		if conn.thingId == "" || id == conn.thingId {
			err := conn.WriteJSON(message{
				Id:          id,
				MessageType: messageType,
				Data:        data,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("websocket closed")
}

func (conn *Connection) onPropertyChanged(msg topic.ThingPropertyChangedMessage) {
	err := conn.write(msg.ThingId, constant.PropertyChanged, map[string]any{msg.PropertyName: msg.Value})
	if err != nil {
		conn.close()
	}
}

func (conn *Connection) onThingRemoved(msg topic.ThingRemovedMessage) {
	err := conn.write(msg.ThingId, constant.ThingRemoved, nil)
	if err != nil {
		conn.close()
	}
}

func (conn *Connection) onThingAdded(msg topic.ThingAddedMessage) {
	err := conn.write(msg.ThingId, constant.ThingAdded, msg.Data)
	if err != nil {
		conn.close()
	}
}

func (conn *Connection) onThingModified(msg topic.ThingModifyMessage) {
	err := conn.write(msg.ThingId, constant.ThingModified, nil)
	if err != nil {
		conn.close()
	}
}

func (conn *Connection) onThingConnected(msg topic.ThingConnectedMessage) {
	err := conn.write(msg.ThingId, constant.Connected, msg.Connected)
	if err != nil {
		conn.close()
	}
}

func (conn *Connection) notify() {
	if conn.thingId == "" {
		ts := conn.container.GetThings()
		copes := make([]*things.Thing, len(ts))
		copy(copes, ts)
		for _, t := range copes {
			err := conn.write(t.GetId(), constant.Connected, t.Connected)
			if err != nil {
				conn.log.Error(err.Error())
				return
			}
			data := make(map[string]any, 0)
			for name, _ := range t.Properties {
				value, err := t.GetPropertyValue(name)
				if err != nil {
					continue
				}
				data[name] = value
			}
			err = conn.write(t.GetId(), constant.PropertyStatus, data)
			if err != nil {
				conn.log.Error(err.Error())
				return
			}
		}
	} else {
		t := conn.container.GetThing(conn.thingId)
		if t != nil {
			err := conn.write(t.GetId(), constant.Connected, t.Connected)
			if err != nil {
				conn.log.Error(err.Error())
				return
			}
			data := make(map[string]any, 0)
			for name, _ := range t.Properties {
				value, err := t.GetPropertyValue(name)
				if err != nil {
					continue
				}
				data[name] = value
			}
			err = conn.write(t.GetId(), constant.PropertyStatus, data)
			if err != nil {
				conn.log.Error(err.Error())
				return
			}
		}

	}
}

func (conn *Connection) handler() {
	defer conn.close()
	conn.log.Infof("新建一个连接：%v", conn.connId)
	_ = conn.container.Subscribe(topic.ThingAdded, conn.onThingAdded)
	_ = conn.container.Subscribe(topic.ThingModify, conn.onThingModified)
	_ = conn.container.Subscribe(topic.ThingConnected, conn.onThingConnected)
	_ = conn.container.Subscribe(topic.ThingPropertyChanged, conn.onPropertyChanged)
	conn.notify()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		conn.handleMessage(data)
	}
}

func (conn *Connection) handleMessage(data []byte) {
	conn.log.Infof("received message:%v", data)
}

func (conn *Connection) close() {
	conn.log.Infof("关闭连接: %v", conn.connId)
	conn.locker.Lock()
	defer conn.locker.Unlock()
	conn.container.Unsubscribe(topic.ThingAdded, conn.onThingAdded)
	conn.container.Unsubscribe(topic.ThingModify, conn.onThingModified)
	conn.container.Unsubscribe(topic.ThingConnected, conn.onThingConnected)
	conn.container.Unsubscribe(topic.ThingPropertyChanged, conn.onPropertyChanged)
	if !conn.isClosed && conn.Conn != nil {
		err := conn.Close()
		if err != nil {
			conn.log.Infof("websocket close error: %s", err.Error())
		}
	}
	conn.isClosed = true
}

func handleWebsocket(container things.Container, log logging.Logger) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		thingId := ws.Params("thingId")
		connectionId = connectionId + 1
		conn := Connection{
			Conn:      ws,
			thingId:   thingId,
			isClosed:  false,
			locker:    sync.Mutex{},
			container: container,
			log:       log,
			connId:    connectionId,
		}
		conn.handler()
	}
}

//func handleWebsocket(container things.Container, log logging.Logger) func(ws *websocket.Conn) {
//	return func(ws *websocket.Conn) {
//		locker := sync.Mutex{}
//		isClosed := false
//		//subscribedEventNames := make(map[string]bool)
//		thingCleanups := make(map[string]func())
//		closeWs := func() {
//			locker.Lock()
//			defer locker.Unlock()
//			if !isClosed {
//				log.Infof("closed websocket connection")
//				ws.Close()
//			}
//			isClosed = true
//			if thingCleanups != nil && len(thingCleanups) > 0 {
//				for id, f := range thingCleanups {
//					log.Infof("thing %s thingCleanup %v", id, f)
//					f()
//				}
//			}
//			log.Infof("websocket close")
//		}
//		defer closeWs()
//		thingId := ws.Params("thingId")
//
//		write := func(id string, messageType string, d ...any) {
//			js := struct {
//				Id          string `json:"id"`
//				MessageType string `json:"messageType"`
//				Data        any    `json:"data,omitempty"`
//			}{
//				id,
//				messageType,
//				func() any {
//					if d == nil {
//						return struct{}{}
//					}
//					return d[0]
//				}(),
//			}
//			log.Infof("write message: %v", js)
//			locker.Lock()
//			defer locker.Unlock()
//			if isClosed {
//				return
//			}
//			if ws.Conn != nil || !isClosed {
//				ws.SetWriteDeadline(time.Now().Add(writeWait))
//				err := ws.WriteJSON(js)
//				if err != nil {
//					log.Error("websocket error : %s", err.Error())
//					return
//				}
//			}
//		}
//
//		//addThing := func(t *things.Thing) {
//		//onThingConnected := func(message topic.ThingConnectedMessage) {
//		//	if message.ThingId != t.GetId() {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.Connected, message.Connected)
//		//}
//		//log.Infof("addThing %s", t.Id)
//		//write(t.GetId(), constant.Connected, t.Connected)
//		//removeConnectedFunc := container.Subscribe(topic.ThingConnected, onThingConnected)
//
//		//onThingRemoved := func(message topic.ThingRemovedMessage) {
//		//	if message.ThingId != t.GetId() {
//		//		return
//		//	}
//		//	//f, ok := thingCleanups[t.GetId()]
//		//	//if ok {
//		//	//	f()
//		//	//}
//		//	delete(thingCleanups, t.GetId())
//		//	write(t.GetId(), constant.ThingRemoved)
//		//}
//		//removeRemovedFunc := container.Subscribe(topic.ThingRemoved, onThingRemoved)
//
//		//onThingModified := func(message topic.ThingModifyMessage) {
//		//	if message.ThingId != t.GetId() {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.ThingModified)
//		//}
//		//removeModifiedFunc := container.Subscribe(topic.ThingModify, onThingModified)
//
//		//onEvent := func(e *events.Event) {
//		//	write(t.GetId(), constant.Event, struct {
//		//		Name  string        `json:"name"`
//		//		Event *events.Event `json:"events"`
//		//	}{Name: e.GetName(), Event: e})
//		//}
//		//removeThingEventStatusFunc := container.Subscribe(topic.ThingEvent, onEvent)
//
//		//onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
//		//	if message.ThingId != t.GetId() {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
//		//}
//		//removePropertyChangedFunc := container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)
//
//		//onActionStatus := func(thingId string, action *actions.ActionDescription) {
//		//	if t.GetId() != thingId {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
//		//}
//		//removeActionStatusFunc := container.Subscribe(topic.ThingActionStatus, onActionStatus)
//
//		//clean := func() {
//		//removeConnectedFunc()
//		//removeRemovedFunc()
//		//removeModifiedFunc()
//		//removePropertyChangedFunc()
//		//removeActionStatusFunc()
//		//removeThingEventStatusFunc()
//		//}
//		//thingCleanups[t.GetId()] = clean
//
//		//}
//
//		handleMessage := func(data []byte) {
//			log.Infof("received message: %v", data)
//		}
//
//		onThingConnected := func(id string) func(message topic.ThingConnectedMessage) {
//			return func(message topic.ThingConnectedMessage) {
//				if message.ThingId != id {
//					return
//				}
//				write(id, constant.Connected, message.Connected)
//			}
//		}
//
//		onThingRemoved := func(id string) func(message topic.ThingRemovedMessage) {
//			return func(message topic.ThingRemovedMessage) {
//				if message.ThingId != id {
//					return
//				}
//				//delete(thingCleanups, id)
//				write(id, constant.ThingRemoved)
//
//			}
//		}
//
//		onThingModified := func(id string) func(message topic.ThingRemovedMessage) {
//			return func(message topic.ThingRemovedMessage) {
//				if message.ThingId != id {
//					return
//				}
//				write(id, constant.ThingModified)
//			}
//		}
//
//		onPropertyChanged := func(id string) func(message topic.ThingPropertyChangedMessage) {
//			return func(message topic.ThingPropertyChangedMessage) {
//				if message.ThingId != id {
//					return
//				}
//				write(id, constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
//			}
//		}
//
//		onActionStatus := func(id string) func(aid string, action *actions.ActionDescription) {
//			return func(aid string, action *actions.ActionDescription) {
//				if aid != id {
//					return
//				}
//				write(id, constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
//			}
//		}
//
//		//onActionStatus := func(thingId string, action *actions.ActionDescription) {
//		//	if t.GetId() != thingId {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
//		//}
//		//removeActionStatusFunc := container.Subscribe(topic.ThingActionStatus, onActionStatus)
//
//		//onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
//		//	if message.ThingId != t.GetId() {
//		//		return
//		//	}
//		//	write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
//		//}
//		//removePropertyChangedFunc := container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)
//
//		onEvent := func(id string) func(e *events.Event) {
//			return func(e *events.Event) {
//				write(id, constant.Event, struct {
//					Name  string        `json:"name"`
//					Event *events.Event `json:"events"`
//				}{Name: e.GetName(), Event: e})
//			}
//		}
//
//		add := func(t1 *things.Thing) {
//			container.Subscribe(topic.ThingConnected, onThingConnected(t1.GetId()))
//			container.Subscribe(topic.ThingRemoved, onThingRemoved(t1.GetId()))
//			container.Subscribe(topic.ThingModify, onThingModified(t1.GetId()))
//			container.Subscribe(topic.ThingEvent, onEvent(t1.GetId()))
//			container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged(t1.GetId()))
//			container.Subscribe(topic.ThingActionStatus, onActionStatus(t1.GetId()))
//			write(t1.GetId(), constant.Connected, t1.Connected)
//			var data = make(map[string]any)
//			for n, _ := range t1.Properties {
//				v, err := t1.GetPropertyValue(n)
//				if err == nil {
//					data[n] = v
//				}
//				write(t1.GetId(), constant.PropertyStatus, data)
//			}
//		}
//
//		onThingAdded := func(message topic.ThingAddedMessage) {
//			locker.Lock()
//			locker.Unlock()
//			if ws.Conn == nil || isClosed {
//				return
//			}
//			err := ws.WriteJSON(map[string]any{
//				"id":          message.ThingId,
//				"messageType": constant.ThingAdded,
//				"data":        struct{}{},
//			})
//			if err != nil {
//				log.Error("websocket send %s message err : %s", constant.ThingAdded, err.Error())
//			}
//			thing := container.GetThing(message.ThingId)
//			if thing != nil {
//				add(thing)
//			}
//		}
//
//		if thingId == "" {
//			ts := container.GetThings()
//			for _, t1 := range ts {
//				add(t1)
//			}
//
//			container.Subscribe(topic.ThingAdded, onThingAdded)
//		} else {
//			t2 := container.GetThing(thingId)
//			if t2 == nil {
//				return
//			}
//			add(t2)
//		}
//
//		ws.SetReadDeadline(time.Now().Add(pongWait))
//		ws.SetPongHandler(func(string) error {
//			log.Infof("ipc server pong handler")
//			ws.SetReadDeadline(time.Now().Add(pongWait))
//			return nil
//		})
//		for {
//			_, message, err := ws.ReadMessage()
//			if err != nil {
//				log.Infof("things websocket read error: %v", err.Error())
//				return
//			}
//			if message == nil {
//				log.Infof("received message: %s", message)
//				handleMessage(message)
//			}
//		}
//	}
//}
