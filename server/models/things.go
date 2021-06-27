package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	AddonManager "github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/things"
	json "github.com/json-iterator/go"
)

type Things struct {
	things.Container
	logger logging.Logger
}

func NewThingsModel(log logging.Logger) *Things {
	instance := &Things{}
	instance.logger = log
	instance.Container = things.NewThingsContainer(things.Options{}, nil, nil, log)
	return instance
}

func (ts *Things) SetPropertyValue(thingId, propName string, value interface{}) (interface{}, error) {
	return nil, nil
}

func (ts *Things) GetPropertyValue(thingId, propName string) (interface{}, error) {
	return nil, nil
}

func (ts *Things) GetPropertiesValue(id string) (map[string]interface{}, error) {
	return nil, nil
}

//func (ts *Things) GetNewThings() []*Thing {
//	connectedDevices := AddonManager.GetDevices()
//	storedThings := ts.GetMapOfThings()
//	var things []*Thing
//	for _, connected := range connectedDevices {
//		for _, storedThing := range storedThings {
//			if connected.GetID() != storedThing.GetID() {
//				data, err := json.MarshalIndent(connected, "", "  ")
//				if err != nil {
//					continue
//				}
//				things = append(things, NewThingFromString(string(data)))
//			}
//		}
//	}
//	return things
//}

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
