package actions

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	json "github.com/json-iterator/go"
)

type ActionInput map[string]any

type Action struct {
	Type        string            `json:"type,omitempty"`
	AtType      string            `json:"@type,omitempty"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Links       []ActionLinksElem `json:"links,omitempty"`
	Forms       []ActionFormsElem `json:"forms,omitempty"`
	Input       ActionInput       `json:"input,omitempty"`
}

func (a *Action) UnmarshalJSON(data []byte) error {
	a.AtType = json.Get(data, "@type").ToString()
	a.Type = json.Get(data, "type").ToString()
	a.Title = json.Get(data, "title").ToString()
	a.Description = json.Get(data, "description").ToString()
	return nil
}

func (a Action) GetType() string {
	return a.Type
}

func (a Action) GetTitle() string {
	return a.Title
}

func (a Action) GetDescription() string {
	return a.Description
}

func (a Action) GetInput() map[string]any {
	return a.Input
}

func (a Action) ToMessage() messages.Action {
	return messages.Action{
		Type:        nil,
		Description: nil,
		Forms:       nil,
		Input:       nil,
		Title:       nil,
	}
}

type ActionDescription struct {
	Id            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Input         map[string]any `json:"input,omitempty"`
	Status        string         `json:"status,omitempty"`
	TimeRequested string         `json:"timeRequested,omitempty"`
	TimeCompleted string         `json:"timeCompleted,omitempty"`
}

type ActionLinksElem struct {
}

type ActionFormsElem struct {
	Op any `json:"op"`
}

func (a ActionDescription) GetName() string {
	return a.Name
}

func (a ActionDescription) GetDescription() any {
	return nil
}
