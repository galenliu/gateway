package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

const (
	ApplicationJson = "application/json"
	LdJSON          = "application/ld+json"
	SenmlJSON       = "application/senml+json"
	CBOR            = "application/cbor"
	SenmlCbor       = "application/senml+cbor"

	XML      = "application/xml"
	SenmlXML = "application/senml+xml"
	EXI      = "application/exi"
)

type DataSchema interface {
	GetType() string
	GetAtType() string
	SetAtType(string)
	SetUnit(string)
	IsReadOnly() bool
	SetType(string)
	SetTitle(s string)
	SetEnum([]interface{})

	MarshalJSON() ([]byte, error)

	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(data []byte) error
}

type dataSchema struct {
	AtType           string        `json:"@type,omitempty"`
	Title            string        `json:"title"`
	Titles           []string      `json:"titles,omitempty"`
	Description      string        `json:"description,omitempty"`
	Descriptions     []string      `json:"descriptions,omitempty"`
	Unit             string        `json:"unit,omitempty"`
	Const            interface{}   `json:"const,omitempty"`
	OneOf            []dataSchema  `json:"oneOf,,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	ReadOnly         bool          `json:"readOnly,omitempty"`
	WriteOnly        bool          `json:"writeOnly,omitempty"`
	Format           string        `json:"format,omitempty"`
	ContentEncoding  string        `json:"contentEncoding,,omitempty"`
	ContentMediaType string        `json:"contentMediaType,,omitempty"`
	Type             string        `json:"type"`
}

func newDataSchemaFromString(description string) *dataSchema {
	data := []byte(description)
	schema := dataSchema{}
	schema.AtType = json.Get(data, "@type").ToString()
	schema.Title = json.Get(data, "title").ToString()
	schema.Titles = json.Get(data, "@type").Keys()
	schema.Description = json.Get(data, "description").ToString()
	schema.Descriptions = json.Get(data, "descriptions").Keys()
	schema.Unit = json.Get(data, "unit").ToString()
	schema.Const = json.Get(data, "const").GetInterface()
	var oneOf []dataSchema
	json.Get(data, "oneOff").ToVal(&oneOf)
	if len(oneOf) > 0 {
		schema.OneOf = oneOf
	}
	var enum []interface{}
	json.Get(data, "enum").ToVal(&enum)
	if len(oneOf) > 0 {
		schema.Enum = enum
	}
	schema.ReadOnly = json.Get(data, "readOnly").ToBool()
	schema.WriteOnly = json.Get(data, "writeOnly").ToBool()
	schema.Format = json.Get(data, "format").ToString()
	schema.ContentEncoding = json.Get(data, "contentEncoding").ToString()
	schema.ContentMediaType = json.Get(data, "contentMediaType").ToString()
	schema.Type = json.Get(data, "type").ToString()
	return &schema
}

func NewDataSchemaFromString(description string) DataSchema {
	data := []byte(description)
	typ := json.Get(data, "type").ToString()
	if typ == "" {
		return nil
	}
	switch typ {
	case hypermedia_controls.TypeNumber:
		return NewNumberSchemaFromString(description)
	case hypermedia_controls.TypeInteger:
		return NewNumberSchemaFromString(description)
	case hypermedia_controls.TypeString:
		return NewStringSchemaFromString(description)
	case hypermedia_controls.TypeArray:
		return NewArraySchemaFromString(description)
	case hypermedia_controls.TypeBoolean:
		return NewBooleanSchemaFromString(description)
	case hypermedia_controls.TypeNull:
		return NewNullSchemaFromString(description)
	case hypermedia_controls.TypeObject:
		return NewObjectSchemaFromString(description)
	default:
		return nil
	}
}

func (d *dataSchema) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&d, "", "  ")
}

func (d *dataSchema) GetType() string {
	return d.Type
}

func (d *dataSchema) GetAtType() string {
	return d.AtType
}

func (d *dataSchema) SetAtType(s string) {
	if s != "" {
		d.AtType = s
	}
}

func (d *dataSchema) SetType(s string) {
	if s != "" {
		if s == hypermedia_controls.TypeNumber || s == hypermedia_controls.TypeInteger || s == hypermedia_controls.TypeObject || s == hypermedia_controls.TypeArray || s == hypermedia_controls.TypeString || s == hypermedia_controls.TypeBoolean || s == hypermedia_controls.TypeNull {
			d.Type = s
		}
	}
}

func (d *dataSchema) GetDescription() string {
	return d.Description
}

func (d *dataSchema) SetTitle(s string) {
	d.Title = s
}

func (d *dataSchema) IsReadOnly() bool {
	return d.ReadOnly
}

func (d *dataSchema) SetUnit(string2 string) {
	d.Unit = string2
}

func (d *dataSchema) SetEnum(e []interface{}) {
	d.Enum = e
}
