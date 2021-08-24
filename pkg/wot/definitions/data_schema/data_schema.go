package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
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
	AtType       string            `json:"@type,omitempty"`
	Title        string            `json:"title,omitempty"`
	Titles       map[string]string `json:"titles,omitempty"`
	Description  string            `json:"description,omitempty"`
	Descriptions map[string]string `json:"descriptions,omitempty"`
	Const        interface{}       `json:"const,omitempty"`
	Default      interface{}       `json:"default,omitempty"`
	Unit         string            `json:"unit,omitempty"`
	OneOf        []DataSchema      `json:"oneOf,,omitempty"`
	Enum         []interface{}     `json:"enum,omitempty"`
	ReadOnly     bool              `json:"readOnly,omitempty"`
	WriteOnly    bool              `json:"writeOnly,omitempty"`
	Format       string            `json:"format,omitempty"`
	Type         string            `json:"type"`
}

func NewDataSchemaFromString(description string) *DataSchema {
	data := []byte(description)
	schema := DataSchema{}
	schema.AtType = controls.JSONGetString(data, "@type", "")
	schema.Title = controls.JSONGetString(data, "title", "")
	schema.Titles = controls.JSONGetMap(data, "titles")
	schema.Description = controls.JSONGetString(data, "description", "")
	schema.Descriptions = controls.JSONGetMap(data, "descriptions")
	schema.Unit = controls.JSONGetString(data, "unit", "")
	schema.Const = json.Get(data, "const").GetInterface()
	schema.Default = json.Get(data, "default").GetInterface()
	var oneOf []DataSchema
	json.Get(data, "oneOff").ToVal(&oneOf)
	if len(oneOf) > 0 {
		schema.OneOf = oneOf
	}
	var enum []interface{}
	json.Get(data, "enum").ToVal(&enum)
	if len(oneOf) > 0 {
		schema.Enum = enum
	}
	schema.ReadOnly = controls.JSONGetBool(data, "readOnly", false)
	schema.WriteOnly = controls.JSONGetBool(data, "writeOnly", false)
	schema.Format = controls.JSONGetString(data, "format", "")
	schema.Type = controls.JSONGetString(data, "type", "")
	if schema.Type == controls.TypeNumber || schema.Type == controls.TypeString || schema.Type == controls.TypeInteger || schema.Type == controls.TypeNull ||
		schema.Type == controls.TypeObject || schema.Type == controls.TypeArray || schema.Type == controls.TypeBoolean {
		return &schema
	}
	return nil
}

//func NewDataSchemaFromString(description string) DataSchema {
//	data := []byte(description)
//	typ := controls.JSONGetString(data, "type", "")
//	if typ == "" {
//		return nil
//	}
//	switch typ {
//	case controls.TypeNumber:
//		return NewNumberSchemaFromString(description)
//	case controls.TypeInteger:
//		return NewNumberSchemaFromString(description)
//	case controls.TypeString:
//		return NewStringSchemaFromString(description)
//	case controls.TypeArray:
//		return NewArraySchemaFromString(description)
//	case controls.TypeBoolean:
//		return NewBooleanSchemaFromString(description)
//	case controls.TypeNull:
//		return NewNullSchemaFromString(description)
//	case controls.TypeObject:
//		return NewObjectSchemaFromString(description)
//	default:
//		return nil
//	}
//}

func (d *DataSchema) GetType() string {
	return d.Type
}

func (d *DataSchema) GetAtType() string {
	return d.AtType
}

func (d *DataSchema) GetDescription() string {
	return d.Description
}

func (d *DataSchema) IsReadOnly() bool {
	return d.ReadOnly
}
