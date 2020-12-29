package addons

import "fmt"

const (
	STRING  = "string"
	BOOLEAN = "boolean"
	INTEGER = "integer"
	NUMBER  = "number"

	UnitHectopascal = "hectopascal"
	UnitKelvin      = "kelvin"

	AlarmProperty                    = "AlarmProperty"
	BarometricPressureProperty       = "BarometricPressureProperty"
	BooleanProperty                  = "BooleanProperty"
	BrightnessProperty               = "BrightnessProperty"
	ColorModeProperty                = "ColorModeProperty"
	ColorProperty                    = "ColorProperty"
	ColorTemperatureProperty         = "ColorTemperatureProperty"
	ConcentrationProperty            = "ConcentrationProperty"
	CurrentProperty                  = "CurrentProperty"
	DensityProperty                  = "DensityProperty"
	FrequencyProperty                = "FrequencyProperty"
	HeatingCoolingProperty           = "HeatingCoolingProperty"
	HumidityProperty                 = "HumidityProperty"
	ImageProperty                    = "ImageProperty"
	InstantaneousPowerFactorProperty = "InstantaneousPowerFactorProperty"
	InstantaneousPowerProperty       = "InstantaneousPowerProperty"
	LeakProperty                     = "LeakProperty"
	LevelProperty                    = "LevelProperty"
	LockedProperty                   = "LockedProperty"
	MotionProperty                   = "MotionProperty"
	OnOffProperty                    = "OnOffProperty"
	OpenProperty                     = "OpenProperty"
	PushedProperty                   = "PushedProperty"
	SmokeProperty                    = "SmokeProperty"
	TargetTemperatureProperty        = "TargetTemperatureProperty"
	TemperatureProperty              = "TemperatureProperty"
	ThermostatModeProperty           = "ThermostatModeProperty"
	VideoProperty                    = "VideoProperty"
	VoltageProperty                  = "VoltageProperty"
)

type IPropertyEvent interface {
	OnPropertyChanged(interface{})
}

type Property struct {
	AtType      string `json:"@type"` //引用的类型
	Type        string `json:"type"`  //数据的格式
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`

	Unit     string `json:"unit,omitempty"` //属性的单位
	ReadOnly bool   `json:"read_only"`
	Visible  bool   `json:"visible"`

	EventEmitter IPropertyEvent

	Minimum interface{} `json:"minimum,omitempty,string"`
	Maximum interface{} `json:"maximum,omitempty,string"`
	Value   interface{} `json:"value"`
	Enum    []interface{}
	Device DeviceProxy
}

func NewStringProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        STRING,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
		Enum:        nil,
	}
	return p
}

func NewBooleanProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        BOOLEAN,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}
	return p
}

func NewNumberProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        NUMBER,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}
	return p
}

func NewIntegerProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        INTEGER,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}
	return p
}

func (prop *Property) setValue(value interface{}) (interface{},error,) {

	switch prop.Type {
	case INTEGER:
		newValue, ok := value.(int64)
		if !ok {
			return prop.Value,fmt.Errorf("value type err")
		}
		prop.setCachedValue(newValue)
	case NUMBER:
		newValue, ok := value.(int64)
		if !ok {
			return prop.Value,fmt.Errorf("value type err")
		}
		prop.setCachedValue(newValue)
	case STRING:
		newValue, ok := value.(string)
		if !ok {
			return prop.Value,fmt.Errorf("value type err")
		}
		prop.setCachedValue(newValue)
	case BOOLEAN:
		newValue, ok := value.(bool)
		if !ok {
			return prop.Value,fmt.Errorf("value type err")
		}
		prop.setCachedValue(newValue)

	}
	return prop.Value,nil
}

func (prop *Property) setCachedValueAndNotify(value interface{}) error {

	prop.EventEmitter.OnPropertyChanged(value)
	return nil
}

func (prop *Property) setCachedValue(value interface{}) {
	prop.Value = value
}

func (prop *Property) doPropertyChanged(p *Property) {
	prop.setCachedValue(p.Value)
}
