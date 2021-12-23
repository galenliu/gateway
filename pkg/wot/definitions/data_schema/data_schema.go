package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
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
	AtType       string                 `json:"@type,omitempty" wot:"optional"`
	Title        string                 `json:"title,omitempty" wot:"optional"`
	Titles       controls.MultiLanguage `json:"titles,omitempty" wot:"optional"`
	Description  string                 `json:"description,omitempty" wot:"optional"`
	Descriptions controls.MultiLanguage `json:"descriptions,omitempty" wot:"optional"`
	Const        any                    `json:"const,omitempty" wot:"optional"`
	Default      any                    `json:"default,omitempty" wot:"optional"`
	Unit         string                 `json:"unit,omitempty" wot:"optional"`
	OneOf        []DataSchema           `json:"oneOf,,omitempty" wot:"optional"`
	Enum         []any                  `json:"enum,omitempty" wot:"optional"`
	ReadOnly     bool                   `json:"readOnly,omitempty" wot:"withDefault"`
	WriteOnly    bool                   `json:"writeOnly,omitempty" wot:"withDefault"`
	Format       string                 `json:"format,omitempty" wot:"optional"`
	Type         string                 `json:"type,,omitempty" wot:"optional"`
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
