package hypermedia_controls

type Form struct {
	Href                URI                          `json:"href" wot:"mandatory"`
	ContentType         string                       `json:"contentType,omitempty" wot:"withDefault"`
	ContentCoding       string                       `json:"contentCoding,omitempty" wot:"optional"`
	Security            ArrayOfString                `json:"security,omitempty" wot:"optional"`
	Scopes              ArrayOfString                `json:"scopes,omitempty" wot:"optional"`
	Response            *ExpectedResponse            `json:"response,omitempty" wot:"optional"`
	AdditionalResponses []AdditionalExpectedResponse `json:"additionalResponses,omitempty" wot:"optional"`
	Subprotocol         string                       `json:"subprotocol,omitempty" wot:"optional"`
	Op                  ArrayOfString                `json:"op,omitempty" wot:"withDefault"`
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
	Op_readproperty            = "readproperty"
	Op_writeproperty           = "writeproperty"
	Op_observeproperty         = "observeproperty"
	Op_unobserveproperty       = "unobserveproperty"
	Op_invokeaction            = "invokeaction"
	Op_queryaction             = "queryaction"
	Op_cancelaction            = "cancelaction"
	Op_subscribeEvent          = "subscribeevent"
	Op_unsubscribeevent        = "unsubscribeevent"
	Op_readallproperties       = "readallproperties"
	Op_writeallproperties      = "writeallproperties"
	Op_readmultipleproperties  = "readmultipleproperties"
	Op_writemultipleproperties = "writemultipleproperties"
	Op_observeallProperties    = "observeallproperties"
	Op_unobserveallproperties  = "unobserveallproperties"
	Op_subscribeallevents      = "subscribeallevents"
	Op_unsubscribeallevents    = "unsubscribeallevents"
	Op_queryallactions         = "queryallactions"
)
