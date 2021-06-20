package data_schema

import (
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
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

type DataSchema struct {
	AtType           string        `json:"@type,omitempty"`
	Title            string        `json:"title"`
	Titles           []string      `json:"titles,omitempty"`
	Description      string        `json:"description,omitempty"`
	Descriptions     []string      `json:"descriptions,omitempty"`
	Unit             string        `json:"unit,omitempty"`
	Const            interface{}   `json:"const,omitempty"`
	OneOf            []DataSchema  `json:"oneOf,,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	ReadOnly         bool          `json:"readOnly,omitempty"`
	WriteOnly        bool          `json:"writeOnly,omitempty"`
	Format           string        `json:"format,omitempty"`
	ContentEncoding  string        `json:"contentEncoding,,omitempty"`
	ContentMediaType string        `json:"contentMediaType,,omitempty"`

	Type string `json:"type"`
}

func NewDataSchemaFromString(data string) DataSchemaInterface {
	typ := json.Get([]byte(data), "type").ToString()
	switch typ {
	case hypermedia_controls.Array:
		return NewArraySchemaFromString(data)
	case hypermedia_controls.Boolean:
		return NewBooleanSchemaFromString(data)
	case hypermedia_controls.Number:
		return NewNumberSchemaFromString(data)
	case hypermedia_controls.Integer:
		return NewIntegerSchemaFromString(data)
	case hypermedia_controls.Object:
		return NewObjectSchemaFromString(data)
	case hypermedia_controls.String:
		return NewStringSchemaFromString(data)
	case hypermedia_controls.Null:
		return NewNullSchemaFromString(data)
	default:
		return nil
	}
}

func (d *DataSchema) GetType() string {
	return d.Type
}

func (d *DataSchema) GetAtType() string {
	return d.AtType
}

func (d *DataSchema) SetAtType(s string) {
	if s != "" {
		d.AtType = s
	}
}

func (d *DataSchema) SetType(s string) {
	if s != "" {
		if s == hypermedia_controls.Number || s == hypermedia_controls.Integer || s == hypermedia_controls.Object || s == hypermedia_controls.Array || s == hypermedia_controls.String || s == hypermedia_controls.Boolean || s == hypermedia_controls.Null {
			d.Type = s
		}
	}
}

func (d *DataSchema) SetTitle(s string) {
	d.Title = s
}

func (d *DataSchema) IsReadOnly() bool {
	return d.ReadOnly
}

func (d *DataSchema) SetUnit(string2 string) {
	d.Unit = string2
}

func (d *DataSchema) SetEnum(e []interface{}) {
	d.Enum = e
}

type DataSchemaInterface interface {
	GetType() string
	GetAtType() string
	SetAtType(string)
	SetUnit(string)
	IsReadOnly() bool
	SetType(string)
	SetTitle(s string)
	SetEnum([]interface{})

	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(data []byte) error
}
