package internal

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
	"time"
)

type DeviceProxy interface {
	HandleGetValue() (interface{}, error)
	NotifyValueChanged(property *Property)
}

type GetValueFunc func() (interface{}, error)
type SetValueFunc func(v interface{}) (interface{}, error)
type NotifyPropertyChanged func(value interface{})

type Property struct {
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	AtType      string        `json:"@type"`
	Unit        string        `json:"unit"`
	Description string        `json:"description"`
	Minimum     int           `json:"minimum"`
	Maximum     int           `json:"maximum"`
	Enum        []interface{} `json:"enum"`
	ReadOnly    bool          `json:"readOnly"`
	MultipleOf  int           `json:"MultipleOf"`
	Forms       []interface{} `json:"forms"`
	Value       interface{}   `json:"value"`

	proxy        DeviceProxy
	response     chan interface{}
	getValueFunc GetValueFunc
	setValueFunc SetValueFunc
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

func NewPropertyFromString(propertyDescr string) *Property {
	data := []byte(propertyDescr)
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
	if p.ReadOnly {
		return nil, fmt.Errorf("read-only property")
	}
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
	hasChanged := p.setCachedValueAndNotify(value)
	if hasChanged {
		return
	}
}

func (p *Property) setCachedValueAndNotify(value interface{}) bool {
	oldValue := p.Value
	p.setCachedValue(value)
	var hasChanged = oldValue != p.Value
	if hasChanged {
		p.proxy.NotifyValueChanged(p)
	}
	return hasChanged
}

func (p *Property) setCachedValue(value interface{}) interface{} {
	if p.Type == TypeBoolean {
		p.Value = !!to.Bool(value)
	} else {
		p.Value = value
	}
	return p.Value
}

func (p *Property) DoPropertyChanged(data []byte) {
	var changed = false
	if t := json.Get(data, "type").ToString(); t != "" && p.Type != t {
		p.Type = t
		changed = true
	}
	if changed {

	}
}
