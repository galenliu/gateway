package hypermedia_controls

type Form struct {
	Href                string               `json:"href,omitempty"`
	Rel                 string               `json:"rel,omitempty"` //w3c定义中无这个，解决和wot兼容
	ContentType         string               `json:"contentType,omitempty"`
	ContentCoding       string               `json:"contentCoding,omitempty"`
	Security            interface{}          `json:"security,omitempty"`
	Scopes              string               `json:"scopes,omitempty"`
	Response            *Response            `json:"response,omitempty"`
	AdditionalResponses *AdditionalResponses `json:"additionalResponses,omitempty"`
	Subprotocol         string               `json:"subprotocol,omitempty"`
	Op                  []string             `json:"op,omitempty"`
}

type Response struct {
	ContentType []string `json:"contentType,omitempty"`
}
type AdditionalResponses struct {
	ContentType []string    `json:"contentType,omitempty"`
	Success     bool        `json:"success,omitempty"`
	Schema      interface{} `json:"schema,omitempty"` //TODO :DataSchemaInterface
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
