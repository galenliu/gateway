package messages

import "fmt"
import "encoding/json"

const (
	TypeNull    = "null"
	TypeBoolean = "boolean"
	TypeObject  = "object"
	TypeArray   = "array"
	TypeNumber  = "number"
	TypeInteger = "integer"
	TypeString  = "string"
)

type DeviceProperties map[string]Property

// PropertyEnumElem The possible values of the property
type PropertyEnumElem any

// PropertyValue The value of the property
type PropertyValue any

// Property Description of the Property
type Property struct {
	// The type of the property
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// The type of the property
	AtType *string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// Description of the property
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Enum corresponds to the JSON schema field "enum".
	Enum []PropertyEnumElem `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Forms corresponds to the JSON schema field "forms".
	//Forms []PropertyFormsElem `json:"forms,omitempty" yaml:"forms,omitempty"`

	// The maximum value of the property
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`

	// The minimum value of the property
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`

	// The precision of the value
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

	// The name of the property
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`

	// If the property is read-only
	ReadOnly *bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`

	// The title of the property
	Title *string `json:"title,omitempty" yaml:"title,omitempty"`

	// The unit of the property
	Unit *string `json:"unit,omitempty" yaml:"unit,omitempty"`

	// The value of the property
	Value PropertyValue `json:"value,omitempty" yaml:"value,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Property) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type in Property: required")
	}
	type Plain Property
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Property(plain)
	return nil
}

type DeviceWithoutIdProperties map[string]any

type FormElementProperty any
