package models

import (
	"gateway/addons"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	Name        string `json:"name"`
	AtType      string `json:"@type"` //引用的类型
	Type        string `json:"type"`  //数据的格式
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Unit     string `json:"unit,omitempty"` //属性的单位
	ReadOnly bool   `json:"read_only"`
	Visible  bool   `json:"visible"`

	Minimum interface{} `json:"minimum,omitempty,string"`
	Maximum interface{} `json:"maximum,omitempty,string"`
	Value   interface{} `json:"value" gorm:"_"`
	Enum    []interface{}

	ThingID int
}

func devPropToThingProp(p *addons.Property) *Property {
	prop := &Property{
		Model:       gorm.Model{},
		Name:        p.Name,
		AtType:      p.AtType,
		Type:        p.Type,
		Title:       p.Title,
		Description: p.Description,
		Unit:        p.Unit,
		ReadOnly:    p.ReadOnly,
		Visible:     p.Visible,
		Minimum:     p.Minimum,
		Maximum:     p.Maximum,
		Value:       p.Value,
		Enum:        p.Enum,
	}
	return prop
}
