package models

import (
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	ActionCompleted = "completed"
	ActionPending   = "pending"
	ActionCreated   = "created"
	ActionError     = "error"
	ActionDeleted   = "deleted"
)

type Action struct {
	ID            string                 `json:"-"`
	ThingId       string                 `json:"-"`
	Name          string                 `json:"name"`
	Input         map[string]interface{} `json:"input,omitempty"`
	Href          string                 `json:"href"`
	TimeRequested string                 `json:"timeRequested,omitempty"`
	TimeCompleted time.Time              `json:"timeCompleted,omitempty"`
	Status        string                 `json:"status"`
	Error         string                 `json:"error,omitempty"`
	bus           ActionsBus
	logger        logging.Logger
}

func NewActionModel(name string, input map[string]interface{}, log logging.Logger) *Action {
	a := &Action{}
	a.logger = log
	a.Input = input
	a.ID = uuid.NewV4().String()
	a.Status = ActionCreated
	a.TimeRequested = time.Stamp
	a.Name = name
	a.Error = ""
	return a
}

func (action *Action) updateStatus(newStatus string) {
	if action.Status == newStatus {
		return
	}
	if newStatus == ActionCompleted {
		action.TimeCompleted = time.Now()
	}
	action.Status = newStatus
	action.bus.Publish(constant.ActionStatus, action)
}

//
//func NewAction(data []byte, thingId string) *Action {
//	a := &Action{}
//	if thingId != "" {
//		a.ThingId = thingId
//	}
//	var kv map[string]interface{}
//	json.Get(data).ToVal(&kv)
//	for k := range kv {
//		if k != "" {
//			a.Name = k
//		}
//	}
//	if a.Name == "" {
//		logging.Error("action name invalid")
//		return nil
//	}
//
//	var input map[string]interface{}
//	json.Get(data, a.Name).Get("input").ToVal(&input)
//	if input != nil {
//		a.Input = input
//	}
//	a.TimeRequested = time.Now().Format("2006-01-02 15:04:05")
//	a.ID = GetUUID()
//	a.Href = fmt.Sprintf("%s/%s/%s", constant.ActionsPath, a.Name, a.ID)
//	if a.ThingId != "" {
//		a.Href = fmt.Sprintf("%s/%s%s", constant.ThingsPath, a.ThingId, a.Href)
//	}
//	a.Status = ActionCreated
//	return a
//}
//

//
//func (action *Action) SetError(msg string) {
//	action.Error = msg
//}
//
//func (action *Action) Subscribe(typ string, f interface{}) {
//	_ = event_bus.Subscribe("Actions."+typ, f)
//}
//
//func (action *Action) Unsubscribe(typ string, f interface{}) {
//	_ = event_bus.Unsubscribe("Actions."+typ, f)
//}
//
//func (action *Action) Publish(typ string, args ...interface{}) {
//	event_bus.Publish("Actions."+typ, args...)
//}
