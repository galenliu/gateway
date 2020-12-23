package models

import (
	"fmt"
	"gateway/addons"
)

const ActionPair = "pair"

var actionId = 0



type Actions struct {
	ActionsList map[int]*Action
}

func NewActions() *Actions {
	return &Actions{
		ActionsList: make(map[int]*Action),
	}
}

func (actions *Actions) AddAction(action *Action) {

	actions.ActionsList[action.ID] = action
	switch action.Name {
	case ActionPair:
		timeout:=action.Input.Get("timeout").ToInt()
		err:=addons.AddNewThing(timeout)
		if err != nil {
			action.Error = err.Error()
		}
	}
}

func (actions *Actions) RemoveAction(actionId int) error{
	action,ok := actions.ActionsList[actionId];if !ok{
		return fmt.Errorf("Invaild action id: %v ",actionId)
	}
	if action.Status == "pending"{

	}
	delete(actions.ActionsList,actionId )
	return nil
}

func generateId() int {
	actionId = actionId + 1
	return actionId
}
