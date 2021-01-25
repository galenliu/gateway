package models

import (
	"addon"
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/util"
	"gateway/plugin"
	"gateway/server/models/thing"
	json "github.com/json-iterator/go"
	"sync"
)

var once sync.Once
var things *Things

type Things struct {
	things                         map[string]*thing.Thing
	onNewThingAddSubscriptionFuncs map[interface{}]func(thing *thing.Thing)
	removeSubscriptionFuncs        []func()
}

func NewThings() *Things {
	once.Do(
		func() {
			things = &Things{}
			things.things = make(map[string]*thing.Thing)
			bus.Subscribe(util.PropertyChanged, things.onPropertyChanged)

		},
	)
	return things
}

func (ts *Things) GetThing(id string) *thing.Thing {
	t, ok := ts.things[id]
	if !ok {
		return nil
	}
	return t
}

func (ts *Things) FindThingProperty(thingId string, propertyName string) *thing.Property {
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

//if models things is null,read new things from the database
func (ts *Things) GetThings() []*thing.Thing {
	var list []*thing.Thing
	if len(ts.things) > 0 {
		for _, t := range ts.things {
			list = append(list, t)
			return list
		}
	}
	return GetThingsFormDataBase()
}

//get things with out database
func (ts *Things) GetNewThings() []*thing.Thing {
	connectedThings := plugin.GetThings()
	storedThings := ts.GetThings()
	var things []*thing.Thing
	var newList []*thing.Thing
	for _, connected := range connectedThings {
		for _, storedThing := range storedThings {
			if connected.ID != storedThing.ID {
				newList = append(newList, connected)
			}
		}
	}
	return things
}

func (ts *Things) HasThing(thingId string) bool {
	_, ok := ts.things[thingId]
	return ok
}



func (ts *Things) CreateThing(id string, description []byte) error {
	var th = thing.NewThing(id, description)
	if th == nil {
		return fmt.Errorf("thing description invaild")
	}
	th.SetConnected(true)
	err := database.CreateThing(th.ID, th.GetDescription())
	return err
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) {
	var th = ts.GetThing(thingId)
	if th == nil {
		return
	}
	prop := th.GetProperty(propName)
	if prop == nil {
		return
	}
	_ = plugin.SetProperValue(thingId, propName, value)
}

func (ts *Things) RemoveThing(thingId string) error {
	//TODO: Delete Thing from database
	t := ts.GetThing(thingId)
	if t == nil {
		return fmt.Errorf("thing not found")
	}
	t.Remove()
	delete(ts.things, thingId)
	return nil
}

func (ts *Things) AddThing(thing *thing.Thing) error {
	//TODO: Delete Thing from database
	return nil
}

func (ts *Things) onPropertyChanged(property *addon.Property) {
	for _, th := range ts.things {
		if th.ID == property.DeviceId {
			for _, prop := range th.Properties {
				if prop.Name == property.Name {
					prop.Value = property.Value
				}
			}
		}
	}
}

func (ts *Things) AddNewThingSubscription(key interface{}, f func(thing *thing.Thing)) func() {
	if ts.onNewThingAddSubscriptionFuncs == nil {
		ts.onNewThingAddSubscriptionFuncs = make(map[interface{}]func(thing *thing.Thing), 2)
	}
	ts.onNewThingAddSubscriptionFuncs[key] = f
	var removeFunc = func() {
		delete(ts.onNewThingAddSubscriptionFuncs, key)
	}
	return removeFunc
}

func GetThingById(thingId string) (thing *thing.Thing, err error) {
	var t string
	t, err = database.QueryValue(thingId)
	err = json.UnmarshalFromString(t, thing)
	return
}

func GetThingsFormDataBase() (list []*thing.Thing) {

	var ts = database.GetThings()
	if ts == nil {
		return nil
	}
	for id, des := range ts {
		t := thing.NewThing(id, []byte(des))
		if t != nil {
			list = append(list, t)
		}
	}
	return
}
