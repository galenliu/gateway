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

type Schema interface {
	GetType() controls.DataSchemaType
	IsReadOnly() bool
}

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

func MarshalSchema(raw json.Any) (s Schema, e error) {
	dataType := raw.Get("type").ToString()
	if dataType == "" {
		return nil, fmt.Errorf("schema type missing")
	}
	switch dataType {
	case controls.TypeInteger:
		var dataSchema IntegerSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeNumber:
		var dataSchema NumberSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeBoolean:
		var dataSchema BooleanSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeArray:
		var dataSchema ArraySchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeObject:
		var dataSchema ObjectSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeNull:
		var dataSchema NullSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	case controls.TypeString:
		var dataSchema StringSchema
		raw.ToVal(&dataSchema)
		if raw.LastError() != nil {
			return nil, raw.LastError()
		}
		return dataSchema, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", dataType)
	}
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
