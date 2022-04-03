package properties

type PropertyType = string
type Type = string
type Unit = string

const (
	Brightness       = "brightness"
	Hue              = "hue"
	ColorTemperature = "ct"
	ColorModel       = "color_mode"

	TypeString  Type = "string"
	TypeBoolean Type = "boolean"
	TypeInteger Type = "integer"
	TypeNumber  Type = "number"

	UnitHectopascal Unit = "hectopascal"
	UnitKelvin      Unit = "kelvin"
	UnitPercentage  Unit = "percentage"
	UnitPercent     Unit = "percent"
	UnitWatt        Unit = "watt"
	UnitHertz       Unit = "hertz"
	UnitVolt        Unit = "volt"

	UnitArcDegrees Unit = "arcdegrees"
	UnitAmpere     Unit = "ampere"
	UnitCelsius    Unit = "celsius"
	UnitLux        Unit = "lux"
	UnitSeconds    Unit = "seconds"
	UnitPPM        Unit = "ppm"

	TypeAlarmProperty                    PropertyType = "AlarmProperty"
	TypeBarometricPressureProperty       PropertyType = "BarometricPressureProperty"
	TypeBooleanProperty                  PropertyType = "BooleanProperty"
	TypeBrightnessProperty               PropertyType = "BrightnessProperty"
	TypeColorModeProperty                PropertyType = "ColorModeProperty"
	TypeColorProperty                    PropertyType = "ColorProperty"
	TypeColorTemperatureProperty         PropertyType = "ColorTemperatureProperty"
	TypeConcentrationProperty            PropertyType = "ConcentrationProperty"
	TypeCurrentProperty                  PropertyType = "CurrentProperty"
	TypeDensityProperty                  PropertyType = "DensityProperty"
	TypeFrequencyProperty                PropertyType = "FrequencyProperty"
	TypeHeatingCoolingProperty           PropertyType = "HeatingCoolingProperty"
	TypeHumidityProperty                 PropertyType = "HumidityProperty"
	TypeImageProperty                    PropertyType = "ImageProperty"
	TypeInstantaneousPowerFactorProperty PropertyType = "InstantaneousPowerFactorProperty"
	TypeInstantaneousPowerProperty       PropertyType = "InstantaneousPowerProperty"
	TypeLeakProperty                     PropertyType = "LeakProperty"
	TypeLevelProperty                    PropertyType = "LevelProperty"
	TypeLockedProperty                   PropertyType = "LockedProperty"
	TypeMotionProperty                   PropertyType = "MotionProperty"
	TypeOnOffProperty                    PropertyType = "OnOffProperty"
	TypeOpenProperty                     PropertyType = "OpenProperty"
	TypePushedProperty                   PropertyType = "PushedProperty"
	TypeSmokeProperty                    PropertyType = "SmokeProperty"
	TypeTargetTemperatureProperty        PropertyType = "TargetTemperatureProperty"
	TypeTemperatureProperty              PropertyType = "TemperatureProperty"
	TypeThermostatModeProperty           PropertyType = "ThermostatModeProperty"
	TypeVideoProperty                    PropertyType = "VideoProperty"
	TypeVoltageProperty                  PropertyType = "VoltageProperty"
)
