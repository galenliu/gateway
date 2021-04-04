package models

import (
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/util"
	"gateway/plugin"
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
			_ = bus.Subscribe(util.ThingAdded, instance.handleNewThing)
		},
	)
	return instance
}

func (ts *Things) GetThing(id string) *Thing {
	t, ok := ts.things["/things/"+id]
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
		ts.things[t.ID] = t
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
	var connectedThings []*Thing
	//connectedThings = new([]*thing.Thing)
	connectedThings = plugin.GetThings()
	//bus.Publish(bus.GetThings, &connectedThings)
	storedThings := ts.GetThings()
	var things []*Thing
	var newList []*Thing
	for _, connected := range connectedThings {
		for _, storedThing := range storedThings {
			if connected.ID != storedThing.ID {
				newList = append(newList, connected)
			}
		}
	}
	return things
}

func (ts *Things) CreateThing(id string, description []byte) (string, error) {
	var th = NewThing(id, description)
	if th == nil {
		return "", fmt.Errorf("thing description invaild")
	}
	err := database.CreateThing(th.ID, th.GetDescription())
	if err != nil {
		return "", err
	}
	ts.things[th.ID] = th
	ts.Publish(util.ThingAdded, th)
	return th.GetDescription(), err
}

func (ts *Things) handleNewThing(data []byte) {
	id := gjson.GetBytes(data, "id").String()
	t := ts.GetThing(id)
	if t == nil {
		return
	}
	t.update(NewThing(id, data))
	t.setConnected(true)
}

func (ts *Things) RemoveThing(thingId string) error {
	//TODO: Delete Thing from database
	//t := ts.GetThing(thingId)
	//if t == nil {
	//	return fmt.Errorf("thing not found")
	//}
	id := "/things/" + thingId
	err := database.RemoveThing(id)
	if err != nil {
		return err
	}
	delete(ts.things, id)
	return nil
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) ([]byte, error) {
	var th = ts.GetThing(thingId)
	if th == nil {
		return nil, fmt.Errorf("thing can not found")
	}
	prop := th.GetProperty(propName)
	if prop == nil {
		return nil, fmt.Errorf("propertyName not found")
	}
	return plugin.SetProperty(thingId, propName, value)

}

func GetThingsFormDataBase() (list []*Thing) {

	var ts = database.GetThings()
	if ts == nil {
		return nil
	}
	for id, des := range ts {
		t := NewThing(id, []byte(des))
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
