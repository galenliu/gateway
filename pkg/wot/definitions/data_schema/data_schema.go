package data_schema

import (
	"fmt"
	"github.com/bytedance/sonic"
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

func MarshalSchema(data []byte) (s Schema, e error) {
	node, err := sonic.Get(data, "type")
	dataType, err := node.String()
	if err != nil {
		return nil, err
	}
	switch dataType {
	case controls.TypeInteger:
		var dataSchema IntegerSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeNumber:
		var dataSchema NumberSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeBoolean:
		var dataSchema BooleanSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeArray:
		var dataSchema ArraySchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeObject:
		var dataSchema ObjectSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeNull:
		var dataSchema NullSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeString:
		var dataSchema StringSchema
		err = sonic.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
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
