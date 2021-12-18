package data_schema

import (
	"fmt"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

const (
	ApplicationJson = "application/json"
	LdJSON          = "application/ld+json"
	SenmlJSON       = "application/senml+json"
	CBOR            = "application/cbor"
	SenmlCbor       = "application/senml+cbor"
	XML             = "application/xml"
	SenmlXML        = "application/senml+xml"
	EXI             = "application/exi"
)

type DataSchema struct {
	AtType       string            `json:"@type,omitempty" wot:"optional"`
	Title        string            `json:"title,omitempty" wot:"optional"`
	Titles       map[string]string `json:"titles,omitempty" wot:"optional"`
	Description  string            `json:"description,omitempty" wot:"optional"`
	Descriptions map[string]string `json:"descriptions,omitempty" wot:"optional"`
	Const        interface{}       `json:"const,omitempty" wot:"optional"`
	Default      interface{}       `json:"default,omitempty" wot:"optional"`
	Unit         string            `json:"unit,omitempty" wot:"optional"`
	OneOf        []DataSchema      `json:"oneOf,,omitempty" wot:"optional"`
	Enum         []any             `json:"enum,omitempty" wot:"optional"`
	ReadOnly     bool              `json:"readOnly,omitempty" wot:"withDefault"`
	WriteOnly    bool              `json:"writeOnly,omitempty" wot:"withDefault"`
	Format       string            `json:"format,omitempty" wot:"optional"`
	Type         string            `json:"type,,omitempty" wot:"optional"`
}

func (schema *DataSchema) UnmarshalJSON(data []byte) error {
	schema.AtType = json.Get(data, "@type").ToString()
	schema.Type = json.Get(data, "type").ToString()
	if schema.Type == "" {
		return fmt.Errorf("type must be set")
	}
	schema.Title = json.Get(data, "title").ToString()
	schema.Titles = controls.JSONGetMap(data, "titles")
	schema.Description = json.Get(data, "description").ToString()
	schema.Descriptions = controls.JSONGetMap(data, "descriptions")
	schema.Unit = json.Get(data, "unit").ToString()
	schema.Const = json.Get(data, "const").GetInterface()
	schema.Default = json.Get(data, "default").GetInterface()

	var oneOf []DataSchema
	if o := json.Get(data, "oneOff"); o.LastError() == nil {
		o.ToVal(&oneOf)
		if len(oneOf) > 0 {
			schema.OneOf = oneOf
		}
	}

	var enum []any
	if e := json.Get(data, "enum"); e.LastError() == nil {
		e.ToVal(&oneOf)
		if len(enum) > 0 {
			schema.Enum = enum
		}
	}

	schema.ReadOnly = controls.JSONGetBool(data, "readOnly", false)
	schema.WriteOnly = controls.JSONGetBool(data, "writeOnly", false)
	schema.Format = json.Get(data, "format").ToString()

	if schema.Type == controls.TypeNumber || schema.Type == controls.TypeString || schema.Type == controls.TypeInteger || schema.Type == controls.TypeNull ||
		schema.Type == controls.TypeObject || schema.Type == controls.TypeArray || schema.Type == controls.TypeBoolean {
		return nil
	}
	return fmt.Errorf("type err")
}

func (schema *DataSchema) GetType() string {
	return schema.Type
}

func (schema *DataSchema) GetAtType() string {
	return schema.AtType
}

func (schema *DataSchema) GetDescription() string {
	return schema.Description
}

func (schema *DataSchema) IsReadOnly() bool {
	return schema.ReadOnly
}
