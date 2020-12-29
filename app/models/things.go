package models

import (
	"fmt"
	"gateway/addons"
	"gateway/event"
	"gateway/pkg/database"
	"github.com/gorilla/websocket"
	"sync"
)

var wsID int64 = 0

var once sync.Once
var things *Things

type Things struct {
	things     map[string]*Thing
	websockets map[int64]*websocket.Conn
	cancelFunc func()
}

func NewThings() *Things {
	once.Do(
		func() {
			things := &Things{
			}
			things.things = make(map[string]*Thing)
			things.websockets = make(map[int64]*websocket.Conn)
			things.cancelFunc = event.ListenDiscoverNewDevice(things.HandleNewThing)
		},

	)
	return things
}

func (ts *Things) GetListThings() []*ThingInfo {
	return nil
}

func (ts *Things) GetThings() map[string]*ThingInfo {

	db, _ := database.GetDB()
	var Things []*ThingInfo
	db.Find(&Things)

	var things = make(map[string]*ThingInfo)
	for _, t := range Things {
		things[t.ID] = t
	}
	return things
}

//用addons manager返回的和数据库中对比，返回新的things
func (ts *Things) GetNewThings() map[string]*ThingInfo {
	connectedThings := deviceToThing(addons.GetThings())
	storedThings := ts.GetThings()
	var thingsMap = make(map[string]*ThingInfo)
	for id, t := range connectedThings {
		if storedThings[id] == nil {
			thingsMap[id] = t
		}
	}
	return thingsMap
}

func (ts *Things) HandleNewThing(device addons.Device) {
	for i, c := range ts.websockets {
		err := c.WriteJSON(device)
		if err != nil {
			delete(ts.websockets, i)
		}
	}
}

func (ts *Things) CreateThing(t Thing) error {
	thing, err := GetThingByIdFormDataBase(t.ID)
	if err != nil {
		return err
	}
	ts.things[thing.ID] = thing
	return nil
}

func (ts *Things) GetThingProperty(thingId, propName string) *Property {
	thing, ok := ts.things[thingId]
	if !ok {
		return nil
	}
	prop, ok := thing.ThingInfo.Properties[propName]
	if !ok {
		return nil
	}
	return prop
}

func (ts *Things) SetThingProperty(thingId, propName string, value interface{}) (interface{}, error) {
	return addons.SetThingProperty(thingId, propName, value)
}

func (ts *Things) RemoveThing(thingId string) error {
	//TODO: Delete Thing from database
	return nil
}

func (ts *Things) RegisterWebsocket(ws *websocket.Conn) {
	wsID = wsID + 1
	ts.websockets[wsID] = ws
	removeFunc := event.ListenDiscoverNewDevice(ts.HandleNewThing)
	removeFunc()
}

func GetThing(id string) (*Thing, error) {
	t, ok := things.things[id]
	if !ok {
		err := fmt.Errorf("thing id: %s invaild", id)
		return nil, err
	}
	return t, nil
}
