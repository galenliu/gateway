package messages

import "fmt"
import "encoding/json"

// Device Description of the device
type Device struct {
	// Context of the device
	Context *string `json:"@context,omitempty" yaml:"@context,omitempty"`

	// Type corresponds to the JSON schema field "@type".
	Type []string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// The base href of the device
	BaseHref *string `json:"baseHref,omitempty" yaml:"baseHref,omitempty"`

	// If credentials are required
	CredentialsRequired *bool `json:"credentialsRequired,omitempty" yaml:"credentialsRequired,omitempty"`

	// Description of the device
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// ID of the device
	Id string `json:"id" yaml:"id"`

	// Links corresponds to the JSON schema field "links".
	Links []Link `json:"links,omitempty" yaml:"links,omitempty"`

	// Pin corresponds to the JSON schema field "pin".
	Pin *Pin `json:"pin,omitempty" yaml:"pin,omitempty"`

	// Properties corresponds to the JSON schema field "properties".
	Properties DeviceProperties `json:"properties,omitempty" yaml:"properties,omitempty"`

	// Events corresponds to the JSON schema field "events".
	Events DeviceEvents `json:"events,omitempty" yaml:"events,omitempty"`

	// Actions corresponds to the JSON schema field "actions".
	Actions DeviceActions `json:"actions,omitempty" yaml:"actions,omitempty"`

	// Title of the Device
	Title *string `json:"title,omitempty" yaml:"title,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Device) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["id"]; !ok || v == nil {
		return fmt.Errorf("field id in Addon_Device: required")
	}
	type Plain Device
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Device(plain)
	return nil
}

// Link A link
type Link struct {
	// The href of the link
	Href string `json:"href" yaml:"href"`

	// The media type of the link
	MediaType *string `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`

	// The type of the relationship
	Rel *string `json:"rel,omitempty" yaml:"rel,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Link) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["href"]; !ok || v == nil {
		return fmt.Errorf("field href in Link: required")
	}
	type Plain Link
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Link(plain)
	return nil
}

// The pin of the device
type Pin struct {
	// The pattern of the pin
	Pattern *string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// If the pin is required
	Required bool `json:"required" yaml:"required"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Pin) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["required"]; !ok || v == nil {
		return fmt.Errorf("field required in Pin: required")
	}
	type Plain Pin
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Pin(plain)
	return nil
}

type Any any

type AnyUri string

// Description of the device
type DeviceWithoutId struct {
	// Context of the device
	Context *string `json:"@context,omitempty" yaml:"@context,omitempty"`

	// Type corresponds to the JSON schema field "@type".
	Type []string `json:"@type,omitempty" yaml:"@type,omitempty"`

	// Actions corresponds to the JSON schema field "actions".
	Actions DeviceWithoutIdActions `json:"actions,omitempty" yaml:"actions,omitempty"`

	// The base href of the device
	BaseHref *string `json:"baseHref,omitempty" yaml:"baseHref,omitempty"`

	// If credentials are required
	CredentialsRequired *bool `json:"credentialsRequired,omitempty" yaml:"credentialsRequired,omitempty"`

	// Description of the device
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Events corresponds to the JSON schema field "events".
	Events DeviceWithoutIdEvents `json:"events,omitempty" yaml:"events,omitempty"`

	// Links corresponds to the JSON schema field "links".
	Links []Link `json:"links,omitempty" yaml:"links,omitempty"`

	// Pin corresponds to the JSON schema field "pin".
	Pin *Pin `json:"pin,omitempty" yaml:"pin,omitempty"`

	// Properties corresponds to the JSON schema field "properties".
	Properties DeviceWithoutIdProperties `json:"properties,omitempty" yaml:"properties,omitempty"`

	// Title of the Device
	Title *string `json:"title,omitempty" yaml:"title,omitempty"`
}

type FormElementResponse struct {
	// ContentType corresponds to the JSON schema field "contentType".
	ContentType *string `json:"contentType,omitempty" yaml:"contentType,omitempty"`
}

type FormElement struct {
	// ContentCoding corresponds to the JSON schema field "contentCoding".
	ContentCoding *string `json:"contentCoding,omitempty" yaml:"contentCoding,omitempty"`

	// ContentType corresponds to the JSON schema field "contentType".
	ContentType *string `json:"contentType,omitempty" yaml:"contentType,omitempty"`

	// Href corresponds to the JSON schema field "href".
	Href AnyUri `json:"href" yaml:"href"`

	// Response corresponds to the JSON schema field "response".
	Response *FormElementResponse `json:"response,omitempty" yaml:"response,omitempty"`

	//// Scopes corresponds to the JSON schema field "scopes".
	//Scopes FormElementScopes `json:"scopes,omitempty" yaml:"scopes,omitempty"`

	//// Security corresponds to the JSON schema field "security".
	//Security FormElementSecurity `json:"security,omitempty" yaml:"security,omitempty"`

	// Subprotocol corresponds to the JSON schema field "subprotocol".
	Subprotocol *string `json:"subprotocol,omitempty" yaml:"subprotocol,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *FormElement) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["href"]; !ok || v == nil {
		return fmt.Errorf("field href in FormElement: required")
	}
	type Plain FormElement
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = FormElement(plain)
	return nil
}

type NotificationLevel int

type Scopes any

type ArrayString []string
