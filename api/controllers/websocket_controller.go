package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/log"
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
	connId    int64
}

func (conn *Connection) write(id string, messageType string, data any) error {
	conn.locker.Lock()
	defer conn.locker.Unlock()
	if conn.Conn != nil && !conn.isClosed {
		if conn.thingId == "" || id == conn.thingId {
			err := conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return err
			}
			err = conn.WriteJSON(message{
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
	return errors.New("websocket closed")
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
	t := conn.container.GetThing(msg.ThingId)
	if t != nil {
		err := conn.write(t.GetId(), constant.Connected, t.Connected)
		if err != nil {
			log.Error(err.Error())
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
			log.Error(err.Error())
			return
		}
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
				log.Error(err.Error())
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
				log.Error(err.Error())
				return
			}
		}
	} else {
		t := conn.container.GetThing(conn.thingId)
		if t != nil {
			err := conn.write(t.GetId(), constant.Connected, t.Connected)
			if err != nil {
				log.Error(err.Error())
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
				log.Error(err.Error())
				return
			}
		}

	}
}

func (conn *Connection) handler() {
	defer conn.close()
	log.Infof("新建一个连接：%v", conn.connId)
	_ = conn.container.Subscribe(topic.ThingAdded, conn.onThingAdded)
	_ = conn.container.Subscribe(topic.ThingModify, conn.onThingModified)
	_ = conn.container.Subscribe(topic.ThingConnected, conn.onThingConnected)
	_ = conn.container.Subscribe(topic.ThingPropertyChanged, conn.onPropertyChanged)
	_ = conn.container.Subscribe(topic.ThingRemoved, conn.onThingRemoved)
	conn.notify()
	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		conn.handleMessage(data)
	}
}

func (conn *Connection) handleMessage(data []byte) {
	log.Infof("received message:%v", data)
}

func (conn *Connection) close() {
	log.Infof("关闭连接: %v", conn.connId)
	conn.container.Unsubscribe(topic.ThingAdded, conn.onThingAdded)
	conn.container.Unsubscribe(topic.ThingModify, conn.onThingModified)
	conn.container.Unsubscribe(topic.ThingConnected, conn.onThingConnected)
	conn.container.Unsubscribe(topic.ThingPropertyChanged, conn.onPropertyChanged)
	conn.container.Unsubscribe(topic.ThingRemoved, conn.onThingRemoved)
	if !conn.isClosed && conn.Conn != nil {
		conn.locker.Lock()
		defer conn.locker.Unlock()
		err := conn.Close()
		if err != nil {
			log.Infof("websocket close error: %s", err.Error())
		}
	}
	conn.isClosed = true
}

func handleWebsocket(container things.Container) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		thingId := ws.Params("thingId")
		connectionId = connectionId + 1
		conn := Connection{
			Conn:      ws,
			thingId:   thingId,
			isClosed:  false,
			locker:    sync.Mutex{},
			container: container,
			connId:    connectionId,
		}
		conn.handler()
	}
}
