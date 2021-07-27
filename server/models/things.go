package models

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/things"
	json "github.com/json-iterator/go"
)

type ThingsModel struct {
	things.Container
	logger logging.Logger
}

func NewThingsModel(things things.Container, log logging.Logger) *ThingsModel {
	instance := &ThingsModel{}
	instance.logger = log
	instance.Container = things
	return instance
}

//func (ts *ThingsModel) GetNewThings() []*Thing {
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
func (ts *ThingsModel) handleNewThing(data []byte) {
	id := json.Get(data, "id").ToString()
	t := ts.GetThing(id)
	if t == nil {
		return
	}
	err := ts.UpdateThing(data)
	if err != nil {
		ts.logger.Error("new thing err : &s", err.Error())
		return
	}
}

func (ts *ThingsModel) handleRemoveThing(thingId string) error {
	err := ts.RemoveThing(thingId)
	if err != nil {
		ts.logger.Error("remove thing err: %s", err.Error())
		return err
	}
	return nil
}
