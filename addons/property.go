package addons

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
	Value   interface{}
	Enum    []interface{}
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

func (prop *Property) setValue(value interface{}) bool {

	var oldValue = prop.Value

	var hasChanged = value != oldValue
	if hasChanged {
		prop.setCachedValueAndNotify(value)
	}
	return hasChanged
}

func (prop *Property) setCachedValueAndNotify(value interface{}) {
	prop.setCachedValue(value)
	prop.EventEmitter.OnPropertyChanged(value)
}

func (prop *Property) setCachedValue(value interface{}) {
	if prop.Type == BOOLEAN {
		prop.Value = value.(bool)
		return
	}
	prop.Value = value
}
