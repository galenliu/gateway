package models

import (
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/util"
	AddonManager "gateway/plugin"
	"github.com/tidwall/gjson"
	"sync"
)

var once sync.Once
var instance *Things

type Things struct {
	things  map[string]*Thing
	Actions *Actions
}

func NewThings() *Things {
	once.Do(
		func() {
			instance = &Things{}
			instance.things = make(map[string]*Thing)
			instance.GetThings()
			instance.Actions = NewActions()
			AddonManager.Subscribe(util.ThingAdded, instance.handleNewThing)
		},
	)
	return instance
}

func (ts *Things) GetThing(id string) *Thing {
	t, ok := ts.things[id]
	if !ok {
		return nil
	}
	return t
}

func (ts *Things) FindThingProperty(thingId string, propertyName string) *Property {
	th := ts.GetThing(thingId)
	if th == nil {
		return nil
	}
	property := th.GetProperty(propertyName)
	if property != nil {
		return property
	}
	return nil
}

//if models instance is null,read new instance from the database
func (ts *Things) GetThings() map[string]*Thing {
	if len(ts.things) > 0 {
		return ts.things
	}
	for _, t := range GetThingsFormDataBase() {
		ts.things[t.GetID()] = t
	}
	return ts.things
}

func (ts *Things) GetListThings() (lt []*Thing) {
	for key, t := range ts.GetThings() {
		t.ID = key
		lt = append(lt, t)
	}
	return
}

//get instance with out database
func (ts *Things) GetNewThings() []*Thing {

	connectedThings := AddonManager.GetDevices()

	storedThings := ts.GetThings()

	var things []*Thing
	for _, connected := range connectedThings {
		for _, storedThing := range storedThings {
			if connected.GetID() != storedThing.GetID() {
				things = append(things, NewThing(connected.GetDescription()))
			}
		}
	}
	return things
}

func (ts *Things) CreateThing(id string, description string) (string, error) {
	var th = NewThing(description)
	th.ID = id
	if th == nil {
		return "", fmt.Errorf("thing description invaild")
	}
	err := database.CreateThing(th.GetID(), th.GetDescription())
	if err != nil {
		return "", err
	}
	ts.things[th.GetID()] = th
	go ts.Publish(util.ThingAdded, th)
	return th.GetDescription(), err
}

func (ts *Things) handleNewThing(data string) {
	id := gjson.Get(data, "id").String()
	t := ts.GetThing(id)
	if t == nil {
		return
	}
	t.update(NewThing(data))
	t.setConnected(true)
}

func (ts *Things) RemoveThing(thingId string) error {
	//TODO: Delete Thing from database
	//t := ts.GetThing(thingId)
	//if t == nil {
	//	return fmt.Errorf("thing not found")
	//}

	err := database.RemoveThing(thingId)
	if err != nil {
		return err
	}
	delete(ts.things, thingId)
	return nil
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) ([]byte, error) {
	var th = ts.GetThing(thingId)
	if th == nil {
		return nil, fmt.Errorf("thing(%s) can not found", thingId)
	}
	prop := th.GetProperty(propName)
	if prop == nil {
		return nil, fmt.Errorf("propertyName not found")
	}
	return AddonManager.SetProperty(thingId, propName, value)

}

func GetThingsFormDataBase() (list []*Thing) {
	var ts = database.GetThings()
	if ts == nil {
		return nil
	}
	for _, des := range ts {
		t := NewThing(des)
		if t != nil {
			t.Connected = false
			list = append(list, t)
		}
	}
	return
}

func (ts *Things) Subscribe(typ string, f interface{}) {
	_ = bus.Subscribe("Things."+typ, f)
}

func (ts *Things) Unsubscribe(typ string, f interface{}) {
	_ = bus.Unsubscribe("Things."+typ, f)
}

func (ts *Things) Publish(typ string, args ...interface{}) {
	bus.Publish("Things."+typ, args...)
}
