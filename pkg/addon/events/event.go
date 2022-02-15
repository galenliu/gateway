package events

import messages "github.com/galenliu/gateway/pkg/ipc_messages"

type EventDescription struct {
}

type Event struct {
	AtType      string           `json:"@type,omitempty"`
	Name        string           `json:"name,omitempty"`
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	Links       []EventLinksElem `json:"links,omitempty"`
	Type        string           `json:"type,omitempty"`
	Unit        string           `json:"unit,omitempty"`
	Minimum     float64          `json:"minimum,omitempty"`
	Maximum     float64          `json:"maximum,omitempty"`
	MultipleOf  float64          `json:"multipleOf,omitempty"`
	Enum        []EventEnumElem  `json:"enum,omitempty"`
}

type EventLinksElem struct {
}

type EventEnumElem struct {
}

func (e Event) GetAtType() string {
	return e.AtType
}

func (e Event) GetName() string {
	return e.Name
}

func (e Event) GetTitle() string {
	return e.Title
}

func (e Event) GetDescription() string {
	return e.Description
}

func (e Event) GetLinks() []EventLinksElem {
	return e.Links
}

func (e Event) GetType() string {
	return e.Type
}

func (e Event) GetEnum() []EventEnumElem {
	return e.Enum
}

func (e Event) GetMinimum() float64 {
	return e.Minimum
}

func (e Event) GetMaximum() float64 {
	return e.Maximum
}

func (e Event) ToMessage() messages.Event {
	return messages.Event{
		Type:        nil,
		Description: nil,
		Enum:        nil,
		Forms:       nil,
		Maximum:     nil,
		Minimum:     nil,
		MultipleOf:  nil,
		Name:        nil,
		Title:       nil,
		Unit:        nil,
	}
}
