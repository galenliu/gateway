package data_schema

import (
	"encoding/json"
	"fmt"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/tidwall/gjson"
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

func UnmarshalSchema(data []byte) (s Schema, e error) {
	dataType := gjson.GetBytes(data, "type").String()
	if dataType == "" {
		return nil, fmt.Errorf("invaild type")
	}
	var err error
	switch dataType {
	case controls.TypeInteger:
		var dataSchema IntegerSchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeNumber:
		var dataSchema NumberSchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeBoolean:
		var dataSchema BooleanSchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeArray:
		var dataSchema ArraySchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeObject:
		var dataSchema ObjectSchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeNull:
		var dataSchema NullSchema
		err = json.Unmarshal(data, &dataSchema)
		if err != nil {
			return nil, err
		}
		return dataSchema, nil
	case controls.TypeString:
		var dataSchema StringSchema
		err = json.Unmarshal(data, &dataSchema)
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
