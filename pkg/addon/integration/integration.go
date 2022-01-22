package integration

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/galenliu/gateway/pkg/constant"
	json "github.com/json-iterator/go"
	"sync"
)

type Thing interface {
	GetId() string
	GetProperty(id string) Property
}

type Property interface {
	OnValueChanged(v any)
}

type Integration struct {
	*proxy.IpcClient
	container  sync.Map
	token      string
	thingsPath string
	dbPath     string
}

func NewIntegration() *Integration {
	return &Integration{}
}

func (c *Integration) AddThing(t Thing) {
	c.container.Store(t.GetId(), t)
}

func (c *Integration) Run() {
	c.IpcClient = proxy.NewClient(c, "9500")

}

func (c *Integration) HandleSetProperty(thingId, propName string, v any) {
	c.IpcClient.Send(Message{
		MessageType: constant.PropertyChanged,
		Data: struct {
			ThingId      string `json:"thingId"`
			PropertyName string `json:"propertyName"`
			Value        any    `json:"value"`
		}{
			ThingId:      thingId,
			PropertyName: propName,
			Value:        v,
		},
	})
}

func (c *Integration) OnMessage(data []byte) {
	mt := json.Get(data, "messageType")
	d := json.Get(data, "data")
	if mt.LastError() != nil || d.LastError() != nil {
		fmt.Println("message error")
		return
	}
	switch mt.ToString() {
	case constant.PropertyChanged:
		thingId := d.Get("deviceId").ToString()
		propName := d.Get("propertyName").ToString()
		value := d.Get("value").GetInterface()
		t := c.GetThing(thingId)
		if t != nil {
			property := t.GetProperty(propName)
			property.OnValueChanged(value)
		}
	case constant.ThingModified:
		thingId := d.Get("thingId").ToString()
		ts := proxy.LoadThings(c.dbPath)
		for n, t := range ts {
			if n == thingId {
				c.HandleThingModified(t)
			}
		}

	}
}

func (c *Integration) LoadThings() map[string]things.Thing {
	return proxy.LoadThings(c.dbPath)
}

func (c *Integration) GetThing(id string) Thing {
	v, ok := c.container.Load(id)
	if ok {
		t, ok := v.(Thing)
		if ok {
			return t
		}
	}
	return nil
}

func (c *Integration) SetProperty() {

}

func (c *Integration) HandleThingModified(thing things.Thing) {
	fmt.Println("HandleThingModified")
}

type Message struct {
	MessageType string `json:"messageType"`
	Data        any    `json:"data"`
}
