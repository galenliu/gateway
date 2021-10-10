package plugin

import (
	"fmt"
	"github.com/galenliu/gateway-addon/properties"
	"github.com/galenliu/gateway-grpc"
	"github.com/xiam/to"

	"math"
	"time"
)

type GetValueFunc func() (interface{}, error)
type SetValueFunc func(v interface{}) (interface{}, error)
type NotifyPropertyChanged func(value interface{})

type Property struct {
	device       *Device
	response     chan interface{}
	SetValueFunc SetValueFunc
	*properties.Property
}

func NewProperty(device *Device, property *properties.Property) *Property {
	p := &Property{}
	p.device = device
	p.Property = property
	return p
}

func NewPropertyFormString(dev *Device, propertyDesc string) *Property {
	p := Property{}
	p.device = dev
	p.Property = properties.NerPropertyFormString(propertyDesc)
	return &p
}

func (p *Property) setValue(value interface{}) error {
	if p.ReadOnly {
		return fmt.Errorf("read-only P")
	}
	var numberValue = to.Float64(value)
	if p.Minimum != nil {
		if to.Float64(p.Minimum) > numberValue {
			return fmt.Errorf("value less than minimum: %s", p.Minimum)
		}
	}
	if p.Maximum != nil {
		if to.Float64(p.Maximum) < numberValue {
			return fmt.Errorf("value greater than minimum: %s", p.Maximum)
		}
	}
	if p.MultipleOf != nil {
		if numberValue/to.Float64(p.MultipleOf)-math.Round(numberValue/to.Float64(p.MultipleOf)) != 0 {
			return fmt.Errorf("value is not a multiple of : %s", p.MultipleOf)
		}
	}
	if len(p.Enum) > 0 {
		for e := range p.Enum {
			if e == value {
				break
			}
			return fmt.Errorf("invalid enum value")
		}
	}
	p.setCachedValueAndNotify(value)
	return nil
}

func (p *Property) setCachedValueAndNotify(value interface{}) bool {
	oldValue := p.Value
	p.SetCachedValue(value)
	var hasChanged = oldValue != p.Value
	return hasChanged
}

func (p *Property) SetCachedValue(value interface{}) {
	if p.Type == rpc.Type_name[int32(rpc.Type_boolean)] {
		p.Value = !!to.Bool(value)
	} else {
		p.Value = value
	}
}

func (p *Property) doPropertyChanged(property *rpc.Property) error {
	err := p.setValue(property.Value)
	if err != nil {
		return err
	}
	select {
	case p.response <- property.Value:
	}
	p.Title = property.Title
	p.Description = property.Description
	return nil
}

func (p *Property) SetValue(value interface{}) (interface{}, error) {
	var data map[string]interface{}
	data["deviceId"] = p.device.ID
	data["propertyName"] = p.Name
	data["propertyValue"] = value
	p.device.adapter.sendMsg(rpc.MessageType_DeviceSetPropertyCommand, data)

	setTimeOut := time.After(time.Duration(1 * time.Second))
	for {
		select {
		case <-setTimeOut:
			return nil, fmt.Errorf("time out")
		case v := <-p.response:
			return v, nil
		}

	}
}
