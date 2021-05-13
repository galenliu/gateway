package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

type Action struct {
	ID            string                 `json:"-"`
	ThingId       string                 `json:"-"`
	Name          string                 `json:"name"`
	Input         map[string]interface{} `json:"input,omitempty"`
	Href          string                 `json:"href"`
	TimeRequested string                 `json:"timeRequested,omitempty"`
	TimeCompleted string                 `json:"timeCompleted,omitempty"`
	Status        string                 `json:"status"`
	Error         string                 `json:"error,omitempty"`
}

func NewAction(data []byte, thingId string) *Action {
	a := &Action{}
	if thingId != "" {
		a.ThingId = thingId
	}
	var kv map[string]interface{}
	json.Get(data).ToVal(&kv)
	for k := range kv {
		if k != "" {
			a.Name = k
		}
	}
	if a.Name == "" {
		log.Error("action name invalid")
		return nil
	}

	var input map[string]interface{}
	json.Get(data, a.Name).Get("input").ToVal(&input)
	if input != nil {
		a.Input = input
	}
	a.TimeRequested = time.Now().Format("2006-01-02 15:04:05")
	a.ID = GetUUID()
	a.Href = fmt.Sprintf("%s/%s/%s", util.ActionsPath, a.Name, a.ID)
	if a.ThingId != "" {
		a.Href = fmt.Sprintf("%s/%s%s", util.ThingsPath, a.ThingId, a.Href)
	}
	a.Status = ActionCreated
	return a
}

func (action *Action) UpdateStatus(newStatus string) {
	if action.Status == newStatus {
		return
	}
	if newStatus == ActionCompleted {
		action.TimeCompleted = time.Now().Format("2006-01-02 15:04:05")
	}
	action.Status = newStatus
	bus.Publish(util.ActionStatus, action)
}

func (action *Action) SetError(msg string) {
	action.Error = msg
}

func (action *Action) Subscribe(typ string, f interface{}) {
	_ = bus.Subscribe("Actions."+typ, f)
}

func (action *Action) Unsubscribe(typ string, f interface{}) {
	_ = bus.Unsubscribe("Actions."+typ, f)
}

func (action *Action) Publish(typ string, args ...interface{}) {
	bus.Publish("Actions."+typ, args...)
}
