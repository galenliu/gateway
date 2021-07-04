package wot_models

import (
	"fmt"
	AddonManager "github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/wot"
	"github.com/xiam/to"
	"sync"
)

const (
	ActionPair   = "pair"
	ActionUnpair = "unpair"

	ActionCompleted = "completed"
	ActionPending   = "pending"
	ActionCreated   = "created"
	ActionError     = "error"
	ActionDeleted   = "deleted"
)

var one sync.Once
var _actions *Actions

type Actions struct {
	things  *wot.Things
	actions map[string]*Action
}

func NewActions() *Actions {
	one.Do(
		func() {
			_actions = &Actions{}
			_actions.actions = make(map[string]*Action)
			_actions.things = wot.NewThings()
		},
	)
	return _actions

}

func (as *Actions) Add(action *Action) error {

	as.actions[action.ID] = action
	if action.ThingId != "" {
		t := as.things.GetThing(action.ThingId)
		if t == nil {
			return fmt.Errorf("invalid thing id: %s", action.ThingId)
		}
		ok := t.AddAction(action)
		if !ok {
			return fmt.Errorf("invalid action name: %s", action.Name)
		}
	}
	action.UpdateStatus(ActionPending)

	switch action.Name {
	case ActionPair:
		timeout := to.Float64(action.Input["timeout"])
		err := AddonManager.AddNewThing(timeout)
		if err != nil {
			action.SetError(err.Error())
			action.UpdateStatus(ActionError)
			return err
		}
		break
	case ActionUnpair:
		//id := json.Get(data, "id").ToInt()
		break
	default:
		delete(as.actions, action.ID)
		return fmt.Errorf("invalid action name: %s", action.Name)
	}
	return nil
}

func (as *Actions) Remove(id string) error {
	a, ok := as.actions[id]
	if !ok {
		return fmt.Errorf("invalid action id: %s", id)
	}
	if a.ThingId != "" {
		t := as.things.GetThing(a.ThingId)
		if t == nil {
			return fmt.Errorf("invalid thing id : %s", a.ThingId)
		}
		err := t.RemoveAction(a)
		if err != nil {
			return err
		}
	} else {
		switch a.Name {
		case ActionPair:
			AddonManager.CancelAddNewThing()
			break
		case ActionUnpair:
			AddonManager.CancelRemoveThing(to.String(a.Input["id"]))
			break
		default:
			return fmt.Errorf("invalid action name: %s", a.Name)
		}
	}
	a.UpdateStatus(ActionDeleted)
	delete(as.actions, a.ID)
	return nil
}

func (as *Actions) GetAction(thingId string, name string) (actions []*Action) {
	for _, a := range as.actions {
		if a.ThingId == thingId {
			if name == "" {
				actions = append(actions, a)
			} else {
				if a.Name == name {
					actions = append(actions, a)
				}
			}
		}
	}
	return
}

func (as *Actions) GetGatewayActions(actionName string) (actions []*Action) {
	for _, a := range as.actions {
		if a.Name == actionName {
			actions = append(actions, a)
		}
	}
	return
}

func (as *Actions) onActionStatus(id string) {

}

func (as *Actions) Subscribe(typ string, f interface{}) {
	_ = event_bus.Subscribe("Actions."+typ, f)
}

func (as *Actions) Unsubscribe(typ string, f interface{}) {
	_ = event_bus.Unsubscribe("Actions."+typ, f)
}

func (as *Actions) Publish(typ string, args ...interface{}) {
	event_bus.Publish("Actions."+typ, args...)
}
