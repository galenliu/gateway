package hypermedia_controls

import json "github.com/json-iterator/go"

type Form struct {
	Href                URI                  `json:"href,omitempty"`
	ContentType         string               `json:"contentType,omitempty"`
	ContentCoding       string               `json:"contentCoding,omitempty"`
	Security            []string             `json:"security,omitempty"`
	Scopes              []string             `json:"scopes,omitempty"`
	Response            *Response            `json:"response,omitempty"`
	AdditionalResponses *AdditionalResponses `json:"additionalResponses,omitempty"`
	SubProtocol         string               `json:"subprotocol,omitempty"`
	Op                  []string             `json:"op,omitempty"`
}

type Response struct {
	ContentType []string `json:"contentType,omitempty"`
}
type AdditionalResponses struct {
	ContentType []string    `json:"contentType,omitempty"`
	Success     bool        `json:"success,omitempty"`
	Schema      interface{} `json:"schema,omitempty"` //TODO :dataSchema
}

func NewFormFormString(description string) *Form {
	data := []byte(description)
	f := &Form{}
	if f.Href = URI(JSONGetString(data, "href", "")); f.Href == "" {
		return nil
	}
	f.ContentType = JSONGetString(data, "contentType", "")
	f.ContentCoding = JSONGetString(data, "contentType", "")

	f.Security = JSONGetArray(data, "security")
	f.Scopes = JSONGetArray(data, "scopes")
	f.SubProtocol = JSONGetString(data, "subprotocol", "")

	var r Response
	json.Get(data, "response").ToVal(&r)
	if &r != nil {
		f.Response = &r
	}

	var ar AdditionalResponses
	json.Get(data, "additionalResponses").ToVal(&r)
	if &ar != nil {
		f.AdditionalResponses = &ar
	}

	f.Op = JSONGetArray(data, "op")

	return f
}

/*
Indicates the semantic intention of performing the operation(s)
described by the form.For example, the Property interaction
allows get and set operations.The protocol binding may contain a form for the get operation
and a different form for the set operation.The op attribute indicates
which form is for which and allows the client to select the correct
form for the operation required. op can be assigned one or
more interaction verb(s) each representing a semantic intention of an operation.
*/

const (
	ReadProperty            = "readProperty"
	WriteProperty           = "writeProperty"
	ObserveProperty         = "observeProperty"
	UnobserveProperty       = "unobserveProperty"
	InvokeAction            = "invokeAction"
	SubscribeEvent          = "subscribeEvent"
	UnsubscribeEvent        = "unsubscribeEvent"
	ReadallProperties       = "readAllProperties"
	WriteAllProperties      = "writeAllProperties"
	ReadMultipleProperties  = "readMultipleProperties"
	WriteMultipleProperties = "writeMultipleProperties"
	ObserveAllProperties    = "observeAllProperties"
	UnobserveAllProperties  = "unobserveAllProperties"
)
