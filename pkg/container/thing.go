package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Thing struct {
	*wot.Thing
	Pin                 *addon.DevicePin `json:"pin,omitempty"`
	CredentialsRequired bool             `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability,omitempty"`
	Connected          bool   `json:"connected,,omitempty"`

	//FloorplanVisibility bool `json:"floorplanVisibility"`
	//FloorplanX          uint `json:"floorplanX"`
	//FloorplanY          uint `json:"floorplanY"`
	//LayoutIndex         uint `json:"layoutIndex"`
	GroupId string `json:"groupId,omitempty"`

	container *ThingsContainer
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(id string, description string) (thing *Thing, err error) {
	if id == "" || description == "" {
		return nil, fmt.Errorf("id or description err")
	}
	data := []byte(description)
	t := Thing{
		Thing:               wot.NewThingFromString(description),
		Pin:                 nil,
		CredentialsRequired: json.Get(data, "credentialsRequired").ToBool(),
		SelectedCapability:  json.Get(data, "selectedCapability").ToString(),
		Connected:           json.Get(data, "connected").ToBool(),
		GroupId:             "",
	}

	if len(t.Type) < 1 || t.Id == "" {
		return nil, fmt.Errorf("@type or id err")
	}

	if t.SelectedCapability == "" {
		t.SelectedCapability = t.Type[0]
	}
	if !util.In(t.SelectedCapability, t.Type) {
		return nil, fmt.Errorf("selectedCapability err")
	}
	return &t, nil
}

func (t *Thing) SetSelectedCapability(selectedCapability string) bool {
	for _, s := range t.Type {
		if s == selectedCapability {
			t.SelectedCapability = selectedCapability
			return true
		}
	}
	return false
}

func (t *Thing) SetTitle(title string) bool {
	if t.Title == title {
		return false
	}
	t.Title = title
	return true
}

func (t *Thing) setConnected(connected bool) {
	if t.Connected == connected {
		return
	}
	t.Connected = connected
	t.container.bus.PublishThingConnected(t.GetId(), connected)
}

func (t *Thing) remove() {
	t.container.bus.PublishThingRemoved(t.GetId())
}
