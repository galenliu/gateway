package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	AddonManager "github.com/galenliu/gateway/plugin"
	json "github.com/json-iterator/go"
)

type Things struct {
	things  map[string]*Thing
	logger  logging.Logger
	Actions *Actions
}

func NewThings(log logging.Logger) *Things {
	instance := &Things{}
	instance.logger = log
	instance.things = make(map[string]*Thing)
	instance.GetMapOfThings()
	return instance
}

func (ts *Things) GetThing(id string) *Thing {
	t, ok := ts.things[id]
	if !ok {
		return nil
	}
	return t
}

// GetMapOfThings if models instance is null,read new instance from the database
func (ts *Things) GetMapOfThings() map[string]*Thing {
	if len(ts.things) > 0 {
		return ts.things
	}
	for _, t := range ts.GetThingsFormDataBase() {
		ts.things[t.ID.GetID()] = t
	}
	return ts.things
}

func (ts *Things) GetListThings() (lt []*Thing) {
	for key, t := range ts.GetMapOfThings() {
		t.ID = key
		lt = append(lt, t)
	}
	return
}

func (ts *Things) GetThingsFormDataBase() []*Thing {
	var things = database.GetThings()
	var list []*Thing
	if ts == nil {
		return nil
	}
	for _, des := range things {
		t := NewThingFromString(des)
		if t != nil {
			t.Connected = false
			list = append(list, t)
		}
	}
	return list
}

func (ts *Things) GetNewThings() []*Thing {
	connectedDevices := AddonManager.GetDevices()
	storedThings := ts.GetMapOfThings()
	var things []*Thing
	for _, connected := range connectedDevices {
		for _, storedThing := range storedThings {
			if connected.GetID() != storedThing.GetID() {
				data, err := json.MarshalIndent(connected, "", "  ")
				if err != nil {
					continue
				}
				things = append(things, NewThingFromString(string(data)))
			}
		}
	}
	return things
}

func (ts *Things) CreateThing(id string, description string) (string, error) {

	th := NewThingFromString(description)
	if &th == nil {
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

//Addon Manager New Device Added
func (ts *Things) handleNewThing(data []byte) {
	id := json.Get(data, "id").ToString()
	t := ts.GetThing(id)
	if t == nil {
		return
	}
	t.updateFromString(string(data))
	if !t.Connected {
		t.setConnected(true)
	}
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

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) (interface{}, error) {
	var th = ts.GetThing(thingId)
	if th == nil {
		return nil, fmt.Errorf("thing(%s) can not found", thingId)
	}
	return AddonManager.SetProperty(thingId, propName, value)

}
