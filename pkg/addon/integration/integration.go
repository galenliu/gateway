package integration

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/galenliu/gateway/pkg/constant"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Thing interface {
	GetId() string
	GetProperty(string2 string) Property
}
type Property interface {
	OnValueChanged(v any)
}

type Integration struct {
	*proxy.IpcClient
	container  sync.Map
	token      string
	thingsPath string
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
	}
}

func (c *Integration) UpdateThings() {
	header := url.Values{}
	header.Add("Authorization", c.token)
	data := header.Encode()
	clint := http.Client{}
	request, _ := http.NewRequest("GET", c.thingsPath, strings.NewReader(data))
	response, err := clint.Do(request)
	var ts []*things.Thing
	d, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(d, &ts)
	if err != nil {
		return
	}
	for _, t := range ts {
		c.container.Store(t.GetId(), t)
	}
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

type Message struct {
	MessageType string `json:"messageType"`
	Data        any    `json:"data"`
}
