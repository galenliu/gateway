package internal

import (
	"fmt"
	json "github.com/json-iterator/go"
	"time"
)

type PropertyHandler interface {
	HandleGetValue() (interface{}, error)
}

type GetValueFunc func() (interface{}, error)
type SetValueFunc func(v interface{}) (interface{}, error)
type NotifyValueChanged func(value interface{})

type Property struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	AtType string `json:"@type"`

	Value interface{} `json:"value"`

	response           chan interface{}
	getValueFunc       GetValueFunc
	setValueFunc       SetValueFunc
	notifyValueChanged NotifyValueChanged
}

func NewProperty(des string, getFunc GetValueFunc) *Property {
	prop := NewPropertyFromString(des)
	prop.getValueFunc = getFunc
	if prop.getValueFunc != nil {
		if v, err := prop.getValueFunc(); err != nil {
			prop.Value = v
		}
	}
	return prop
}

func NewPropertyFromString(des string) *Property {
	data := []byte(des)
	p := Property{}
	p.Name = json.Get(data, "name").ToString()
	p.Type = json.Get(data, "type").ToString()
	return nil
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetType() string {
	return p.Type
}

func (p *Property) GetAtType() string {
	return p.AtType
}

func (p *Property) GetValue() (interface{}, error) {
	return p.getValueFunc()
}

func (p *Property) SetValue(value interface{}, await ...int) (interface{}, error) {
	var i = 1000
	if len(await) > 0 {
		i = await[0]
	}
	v, err := p.setValueFunc(value)
	if err != nil {
		p.OnValueChanged(v)
	}
	timeout := time.After(time.Duration(i) * time.Microsecond)
	for {
		select {
		case value := <-p.response:
			return value, nil
		case <-timeout:
			return nil, fmt.Errorf("set property value timeout")
		}
	}
}

func (p *Property) OnValueChanged(value interface{}) {
	err := p.setLocalValue(value)
	if err != nil {
		return
	}

	p.notifyValueChanged(value)
}

func (p *Property) setLocalValue(value interface{}) error {
	p.Value = value
	return nil
}
