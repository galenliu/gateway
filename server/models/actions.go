package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/plugin"
	"sync"
)

const (
	ActionPair    = "pair"
	ActionUnpair  = "unpair"
	ActionPending = "pending"
)

var one sync.Once
var instanceActions *Actions

type Actions struct {
	things *Things
	List   map[string]*Action
}

func NewActions() *Actions {
	one.Do(
		func() {

			instanceActions = &Actions{}
			instanceActions.List = make(map[string]*Action)

		},
	)
	return instanceActions

}

func (actions *Actions) Add(action *Action) {

	actions.List[action.ID] = action
	if action.ThingId != "" {
		t := actions.things.GetThing(action.ThingId)
		if t == nil {
			action.Error = "can not find thing"
			return
		}
	}

	switch action.Name {
	case ActionPair:
		timeout, ok := action.Input["timeout"].(float64)
		if !ok {
			return
		}
		err := plugin.AddNewThing(timeout)
		if err != nil {
			action.Error = err.Error()
		}
		break
	case ActionUnpair:
		//id := json.Get(data, "id").ToInt()
		break
	default:
		delete(actions.List, action.ID)
	}

}

func (actions *Actions) Remove(actionId string) error {
	action, ok := actions.List[actionId]
	if !ok {
		return fmt.Errorf("Invaild actions id: %v ", actionId)
	}
	if action.Status == "pending" {
		if action.ThingId != "" {
			if t := actions.things.GetThing(action.ThingId); t != nil {
				if !t.RemoveAction(action) {
					return fmt.Errorf(fmt.Sprintf("Invaild action name : %s", action.Name))
				}
			}
		}

	} else {
		switch action.Name {
		case ActionPair:
			plugin.CancelAddNewThing()
			break
		case ActionUnpair:
			plugin.CancelRemoveThing(action.Input["id"].(string))
			break
		default:
			return fmt.Errorf("Invaild action name:" + action.Name)
		}
	}
	action.UpdateStatus("deleted")
	delete(actions.List, actionId)
	return nil
}

func (actions *Actions) Subscribe(typ string, f interface{}) {
	_ = bus.Subscribe("Actions."+typ, f)
}

func (actions *Actions) Unsubscribe(typ string, f interface{}) {
	_ = bus.Unsubscribe("Actions."+typ, f)
}

func (actions *Actions) Publish(typ string, args ...interface{}) {
	bus.Publish("Actions."+typ, args...)
}
