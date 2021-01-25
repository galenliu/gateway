package models

import (
	"fmt"
	"gateway/plugin"
	"gateway/server/models/thing"
)

const (
	ActionPair    = "pair"
	ActionUnpair  = "unpair"
	ActionPending = "pending"
)

var actionId uint = 0

type Actions struct {
	things *Things
	List map[uint]*thing.Action
}

func NewActions() *Actions {

	return &Actions{
		List: make(map[uint]*thing.Action),
	}
}

func (actions *Actions) AddAction(action *thing.Action) {
	action.ID = actionId
	actions.List[action.ID] = action

	if action.ThingID != "" {
		t := actions.things.GetThing(action.ThingID)
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

func (actions *Actions) RemoveAction(actionId uint) error {
	action, ok := actions.List[actionId]
	if !ok {
		return fmt.Errorf("Invaild actions id: %v ", actionId)
	}
	if action.Status == "pending" {

	}
	delete(actions.List, actionId)
	return nil
}

func generateId() uint {
	actionId = actionId + 1
	return actionId
}
