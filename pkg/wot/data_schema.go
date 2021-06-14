package wot

import (
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

func NewDataSchemaFromString(data string) IDataSchema {
	typ := json.Get([]byte(data), "type").ToString()
	switch typ {
	case Array:
		return NewArraySchemaFromString(data)
	case Boolean:
		return NewBooleanSchemaFromString(data)
	case Number:
		return NewNumberSchemaFromString(data)
	case Integer:
		return NewIntegerSchemaFromString(data)
	case Object:
		return NewObjectSchemaFromString(data)
	case String:
		return NewStringSchemaFromString(data)
	case Null:
		return NewNullSchemaFromString(data)
	default:
		return nil
	}
}

func (d *DataSchema) GetType() string {
	return d.Type
}

func (d *DataSchema) SetAtType(s string) {
	if s != "" {
		d.AtType = s
	}
}

func (d *DataSchema) SetType(s string) {
	if s != "" {
		if s == Number || s == Integer || s == Object || s == Array || s == String || s == Boolean || s == Null {
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

type IDataSchema interface {
	GetType() string
	SetAtType(string)
	SetUnit(string)
	IsReadOnly() bool
	SetType(string)
	SetTitle(s string)
	SetEnum([]interface{})

	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(data []byte) error
}
