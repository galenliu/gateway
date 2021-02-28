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

type Actions struct {
	things *Things
	List   map[string]*thing.Action
}

func NewActions() *Actions {
	return &Actions{
		List: make(map[string]*thing.Action),
	}
}

func (actions *Actions) Add(action *thing.Action) {

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

func (actions *Actions) Remove(actionId string) error {
	action, ok := actions.List[actionId]
	if !ok {
		return fmt.Errorf("Invaild actions id: %v ", actionId)
	}
	if action.Status == "pending" {
		if action.ThingID != "" {
			if t := actions.things.GetThing(action.ThingID); t != nil {
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
