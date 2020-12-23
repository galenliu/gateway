package models

import (
	"fmt"
	"gateway/addons"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
)

type IThings interface {
	GetProperty(thingId, propName string) interface{}
	SetProperty(thingId, propName string, value interface{})
	RemoveThing(thingId string)
	GetInstallAddons() map[string]interface{}
}

var wsID int64 = 0

type Things struct {
	things     map[string]*thing
	Manager    IThings
	websockets map[int64]*websocket.Conn
}

func NewThings() *Things {
	ts := &Things{
	}
	ts.things = make(map[string]*thing)
	ts.websockets = make(map[int64]*websocket.Conn)
	return ts
}

func (ts *Things) GetListThings() []*Thing {
	return nil
}

func (ts *Things) GetThings() map[string]*Thing {

	db := database.GetDB()
	var Things []*Thing
	db.Find(&Things)

	var things = make(map[string]*Thing)
	for _, t := range Things {
		things[t.ID] = t
	}
	return things
}

func (ts *Things) GetThing(id string) (*thing, error) {
	t, ok := ts.things[id]
	if !ok {
		err := fmt.Errorf("thing id: %v invaild", id)
		return nil, err
	}
	return t, nil
}

//用addons manager返回的和数据库中对比，返回新的things
func (ts *Things) GetNewThings() map[string]*Thing {
	connectedThings := deviceToThing(addons.GetThings())
	storedThings := ts.GetThings()
	var thingsMap = make(map[string]*Thing)
	for id, t := range connectedThings {
		if storedThings[id] == nil {
			thingsMap[id] = t
		}
	}
	return thingsMap
}

func (ts *Things) CreateThing(t Thing) error {
	db := database.GetDB()
	err := db.AutoMigrate(&Thing{})
	err = db.AutoMigrate(&Property{})
	if err != nil {
		return err
	}
	thing := NewThing(t)
	thing.connected = true
	db.First(t)
	ts.things[t.ID] = thing
	return nil
}

func (ts *Things) SetThing(t thing) error {
	return nil
}

func (ts *Things) ToJson() string {
	data, err := json.MarshalToString(ts.things)
	if err != nil {
		log.Error("things marshal err")
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

func (ts *Things) RegisterWebsocket(ws *websocket.Conn) func() {
	wsID = wsID + 1
	ts.websockets[wsID] = ws
	return func() { delete(ts.websockets, wsID) }
}
