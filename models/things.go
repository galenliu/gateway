package models

import (
	"gateway"
	"gateway/plugin"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var log *zap.Logger

type Things struct {
	things       map[string]*Thing
	AddonManager *plugin.AddonsManager
}

func NewThings(manager *plugin.AddonsManager) *Things {
	log = gateway.GetLogger()
	return &Things{
		AddonManager: manager,
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
	return ts.AddonManager.GetProperty(thingId, propName)
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) {
	ts.AddonManager.SetProperty(thingId, propName, value)
}
