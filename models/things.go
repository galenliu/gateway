package models

import (
	json "github.com/json-iterator/go"
	logger "go.uber.org/zap"
)

var log *logger.Logger

type IAddonManager interface {
	GetProperty(thingId, propName string) interface{}
	SetProperty(thingId, propName string, value interface{})
	RemoveThing(thingId string)
	GetInstallAddons() map[string]interface{}
}

type Things struct {
	things  map[string]*Thing
	Manager IAddonManager
}

func NewThings(manager IAddonManager) *Things {

	return &Things{
		Manager: manager,
	}
}

func (ts *Things) GetListThings() []*Thing {
	return nil
}

func (ts *Things) GetThings() map[string]*Thing {
	if len(ts.things) > 0 {
		return ts.things
	}
	return nil
}

func (ts *Things) GetThing(id string) *Thing {
	t := ts.things[id]
	return t
}

func (ts *Things) ToJson() string {
	data, err := json.MarshalToString(ts.things)
	if err != nil {
		log.Warn("things marshal err")
	}
	return data
}

func (ts *Things) GetThingProperty(thingId, propName string) interface{} {
	return ts.Manager.GetProperty(thingId, propName)
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) {
	ts.Manager.SetProperty(thingId, propName, value)
}

func (ts *Things) RemoveThing(thingId string) error {
	//TODO: Delete thing from database
	return nil
}
