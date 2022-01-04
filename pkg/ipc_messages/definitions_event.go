package messages

import json "github.com/json-iterator/go"

type DeviceEvents map[string]Event

// EventEnumElem The possible values of the property
type EventEnumElem any

// Event Description of the events
type Event struct {
	// The type of the events
	Type *string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// Description of the events
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Enum corresponds to the JSON schema field "enum".
	Enum []EventEnumElem `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Forms corresponds to the JSON schema field "forms".
	Forms []EventFormsElem `json:"forms,omitempty" yaml:"forms,omitempty"`

	// The maximum value of the events
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`

	// The minimum value of the events
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`

	// The precision of the value
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

	// The name of the events
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`

	// The title of the events
	Title *string `json:"title,omitempty" yaml:"title,omitempty"`

	// The unit of the events
	Unit *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type DeviceWithoutIdEvents map[string]any

type EventFormsElem struct {
	FormElement
	Op EventFormsElemOp `json:"op"`
}

const (
	EventFormsElemOpSubscribeEvent   = "subscribeEvent"
	EventFormsElemOpUnsubscribeEvent = "unsubscribeEvent"
)

type EventFormsElemOp struct {
	*string
	*ArrayString
}

func (a *EventFormsElemOp) UnmarshalJSON(data []byte) error {
	var raw []string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		a.ArrayString = (*ArrayString)(&raw)
		return nil
	} else {
		str := string(data)
		a.string = &str
	}
	return nil
}

func (a *EventFormsElemOp) MarshalJSON() ([]byte, error) {
	if a.ArrayString != nil {
		bt, err := json.Marshal(*a.ArrayString)
		if err != nil {
			return bt, nil
		}
		return nil, err
	}
	if a.string != nil {
		bt, err := json.Marshal(*a.string)
		if err != nil {
			return bt, nil
		}
		return nil, err
	}
	return nil, nil
}
