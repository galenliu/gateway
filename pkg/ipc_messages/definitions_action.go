package messages

import "encoding/json"

type DeviceActions map[string]Action

// Action Description of the Action
type Action struct {
	// The type of the actions
	Type *string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// Description of the Action
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Forms corresponds to the JSON schema field "forms".
	Forms []ActionFormsElem `json:"forms,omitempty" yaml:"forms,omitempty"`

	// Input corresponds to the JSON schema field "input".
	Input ActionInput `json:"input,omitempty" yaml:"input,omitempty"`

	// The title of the Action
	Title *string `json:"title,omitempty" yaml:"title,omitempty"`
}

const ActionFromOpInvokeAction = "invokeaction"

type ActionFormsElem struct {
	FormElement
	Op ActionFormsElemOp `json:"op"`
}

type ActionFormsElemOp struct {
	*string
	*ArrayString
}

func (a *ActionFormsElemOp) UnmarshalJSON(data []byte) error {
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

type ActionInput any

// ActionInputPropertyEnumElem The possible values of the input property
type ActionInputPropertyEnumElem any

// ActionInputProperty An actions input property
type ActionInputProperty struct {
	// The type of the input property
	Type *string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// Enum corresponds to the JSON schema field "enum".
	Enum []ActionInputPropertyEnumElem `json:"enum,omitempty" yaml:"enum,omitempty"`

	// The maximum value of the input property
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`

	// The minimum value of the input property
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`

	// The precision of the value
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

	// The unit of the input property
	Unit *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type DeviceWithoutIdActions map[string]any
