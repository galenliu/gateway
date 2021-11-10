package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Thing struct {
	*wot.Thing

	bus EventBus

	//The configuration  of the addon
	Pin                 *PIN `json:"pin,omitempty"`
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconHref           string `json:"iconHref,omitempty"`

	FloorplanVisibility bool `json:"floorplanVisibility"`
	FloorplanX          uint `json:"floorplanX"`
	FloorplanY          uint `json:"floorplanY"`
	LayoutIndex         uint `json:"layoutIndex"`

	Security            string             `json:"security"`
	SecurityDefinitions SecurityDefinition `json:"securityDefinitions"`
	GroupId             string             `json:"group_id"`
}

func (t Thing) AddConnectedSubscription(f func(bool)) func() {
	topic := t.Id.GetId() + "." + constant.CONNECTED
	t.bus.bus.Subscribe(topic, f)
	return func() {
		t.bus.bus.Unsubscribe(topic, f)
	}
}

func (t Thing) AddSubscription(topi string, f interface{}) func() {
	topic := t.Id.GetId() + "." + topi
	t.bus.bus.Subscribe(topic, f)
	return func() {
		t.bus.bus.Unsubscribe(topic, f)
	}
}

func (t Thing) setConnected(connected bool) {
	t.Connected = connected
	t.bus.PublishConnected(connected)
}

func (t Thing) remove() {
	t.bus.PublishRemoved()
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(id string, description string) (thing *Thing, err error) {
	if id == "" || description == "" {
		return nil, fmt.Errorf("id or description err")
	}
	data := []byte(description)
	var p PIN
	json.Get(data, "pin").ToVal(&p)

	t := Thing{
		Thing:               wot.NewThingFromString(description),
		Pin:                 &p,
		CredentialsRequired: json.Get(data, "credentialsRequired").ToBool(),
		SelectedCapability:  json.Get(data, "selectedCapability").ToString(),
		Connected:           json.Get(data, "connected").ToBool(),
		IconHref:            json.Get(data, "iconHref").ToString(),
		FloorplanVisibility: false,
		FloorplanX:          0,
		FloorplanY:          0,
		LayoutIndex:         0,
		Security:            "",
		SecurityDefinitions: SecurityDefinition{},
		GroupId:             "",
	}

	if len(t.AtType) < 1 || t.Id == "" {
		return nil, fmt.Errorf("@type or id err")
	}

	if t.SelectedCapability == "" {
		t.SelectedCapability = t.AtType[0]
	}

	if !util.In(t.SelectedCapability, t.AtType) {
		return nil, fmt.Errorf("selectedCapability err")
	}
	return &t, nil
}
