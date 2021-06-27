package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Thing struct {
	*core.Thing
	//The configuration  of the device
	Pin                 *PIN `json:"pin,omitempty"`
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(description string) (*Thing, error) {
	data := []byte(description)
	thing := Thing{}
	t, err := core.NewThingFromString(description)
	if err != nil {
		return nil, err
	}
	thing.Thing = t
	var pin PIN
	json.Get(data, "pin").ToVal(&pin)
	if &pin != nil {
		thing.Pin = &pin
	}
	if c := json.Get(data, "credentialsRequired"); c.ValueType() == json.BoolValue {
		thing.CredentialsRequired = c.ToBool()
	}
	if connected := json.Get(data, "credentialsRequired"); connected.ValueType() == json.BoolValue {
		thing.CredentialsRequired = connected.ToBool()
	}

	if i := json.Get(data, "iconData"); i.ValueType() == json.StringValue {
		thing.IconData = i.ToString()
	}

	return &thing, nil
}

func (t *Thing) setSelectedCapability(s string) {
	if t.SelectedCapability == s {
		return
	}
	for _, typ := range t.AtType {
		if s == typ {
			t.SelectedCapability = s
			err := t.save()
			if err != nil {
				return
			}

		}
	}

}

func (t *Thing) GetSelectedCapability() string {
	return t.SelectedCapability
}

func (t *Thing) SetTitle(title string) string {
	if t.Title != title {
		t.Title = title
		err := t.save()
		if err != nil {
			logging.Info(err.Error())
		}

	}
	return t.GetDescription()
}

func (t *Thing) GetTitle() string {
	return t.Title
}

func (t *Thing) setConnected(connected bool) {
	err := t.save()
	if err != nil {
		logging.Info(err.Error())
	}
	t.Connected = connected
}

func (t *Thing) isConnected() bool {
	return t.Connected
}

func (t *Thing) AddAction(a *Action) bool {
	_, ok := t.Actions[a.Name]
	return ok
}

func (t *Thing) RemoveAction(a *Action) error {
	_, ok := t.Actions[a.Name]
	if !ok {
		return fmt.Errorf("invalid action name :%s", a.Name)
	}
	return nil
}

func (t *Thing) GetDescription() string {
	s, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return ""
	}
	return string(s)
}

func (t *Thing) updateFromString(data string) {

}

func (t *Thing) save() error {
	err := database.UpdateThing(t.GetID(), t.GetDescription())
	if err != nil {
		return err
	}
	return nil
}

//func (t *Thing) Subscribe(typ string, f interface{}) {
//	go func() {
//		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
//		if err != nil {
//			logging.Error(err.Error())
//		}
//	}()
//}
//
//func (t *Thing) Unsubscribe(typ string, f interface{}) {
//	go func() {
//		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
//		if err != nil {
//			logging.Error(err.Error())
//		}
//	}()
//}
//
//func (t *Thing) Publish(typ string, args ...interface{}) {
//	go event_bus.Publish(t.GetID()+"."+typ, args...)
//}
