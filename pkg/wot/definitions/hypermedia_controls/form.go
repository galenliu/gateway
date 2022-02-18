package hypermedia_controls

type OP = string
type ContentType = string

type Form struct {
	Href                URI                          `json:"href" wot:"mandatory"`
	ContentType         ContentType                  `json:"contentType,omitempty" wot:"withDefault"`
	ContentCoding       string                       `json:"contentCoding,omitempty" wot:"optional"`
	Security            ArrayOrString                `json:"security,omitempty" wot:"optional"`
	Scopes              ArrayOrString                `json:"scopes,omitempty" wot:"optional"`
	Response            *ExpectedResponse            `json:"response,omitempty" wot:"optional"`
	AdditionalResponses []AdditionalExpectedResponse `json:"additionalResponses,omitempty" wot:"optional"`
	Subprotocol         string                       `json:"subprotocol,omitempty" wot:"optional"`
	Op                  []OP                         `json:"op,omitempty" wot:"withDefault"`
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

func NewOpArray(args ...OP) []OP {
	var arr = make([]OP, 0)
	for _, a := range args {
		arr = append(arr, a)
	}
	return arr
}

const (
	Readproperty            OP = "readproperty"
	Writeproperty           OP = "writeproperty"
	Observeproperty         OP = "observeproperty"
	Unobserveproperty       OP = "unobserveproperty"
	Invokeaction            OP = "invokeaction"
	Queryaction             OP = "queryaction"
	Cancelaction            OP = "cancelaction"
	SubscribeEvent          OP = "subscribeevent"
	Unsubscribeevent        OP = "unsubscribeevent"
	Readallproperties       OP = "readallproperties"
	Writeallproperties      OP = "writeallproperties"
	Readmultipleproperties  OP = "readmultipleproperties"
	Writemultipleproperties OP = "writemultipleproperties"
	ObserveallProperties    OP = "observeallproperties"
	Unobserveallproperties  OP = "unobserveallproperties"
	Subscribeallevents      OP = "subscribeallevents"
	Unsubscribeallevents    OP = "unsubscribeallevents"
	Queryallactions         OP = "queryallactions"

	JSON   ContentType = "application/json"
	LdJSON ContentType = "application/ld+json"
	xml    ContentType = "application/xml"
)
