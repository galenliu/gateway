package models

import (
	"fmt"
	"gateway/addons"
	json "github.com/json-iterator/go"
)

const (
	ActionPair    = "pair"
	ActionUnpair  = "unpair"
	ActionPending = "pending"
)

var actionId = 0

type Actions struct {
	List map[int]*Action
}

func NewActions() *Actions {

	return &Actions{
		List: make(map[int]*Action),
	}
}

func (actions *Actions) AddAction(action *Action) {
	action.ID = actionId
	actions.List[action.ID] = action

	if action.ThingID != "" {
		thing, err := GetThing(action.ThingID)
		if err != nil {
			action.Error = "can not find thing"
			return
		}
		err = thing.AddAction(action)
		if err != nil {
			delete(actions.List, actionId)
		}
	}
	action.updateStatus(ActionPending)

	data, ok := action.Input.([]byte)
	if !ok {
		action.Error = fmt.Sprintf("acton input err")
		action.updateStatus(Error)
		return
	}

	switch action.Name {
	case ActionPair:
		timeout := json.Get(data, "timeout").ToInt()
		err := addons.AddNewThing(timeout)
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

func (actions *Actions) RemoveAction(actionId int) error {
	action, ok := actions.List[actionId]
	if !ok {
		return fmt.Errorf("Invaild action id: %v ", actionId)
	}
	if action.Status == "pending" {

	}
	delete(actions.List, actionId)
	return nil
}

func generateId() int {
	actionId = actionId + 1
	return actionId
}
