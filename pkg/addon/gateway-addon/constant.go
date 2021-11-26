package addon

const (
	TypeString  = "string"
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"

	UnitHectopascal = "hectopascal"
	UnitKelvin      = "kelvin"
	UnitPercentage  = "percentage"
	UnitArcDegrees  = "arcdegrees"
	UnitCelsius     = "celsius"
	UnitLux         = "lux"
	UnitSeconds     = "seconds"
	UnitPPM         = "ppm"

	AlarmProperty                    = "AlarmProperty"
	BarometricPressureProperty       = "BarometricPressureProperty"
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

	Alarm                    = "Alarm"
	AirQualitySensor         = "AirQualitySensor"
	BarometricPressureSensor = "BarometricPressureSensor"
	BinarySensor             = "BinarySensor"
	Camera                   = "Camera"
	ColorControl             = "ColorControl"
	ColorSensor              = "ColorSensor"
	DoorSensor               = "DoorSensor"
	EnergyMonitor            = "EnergyMonitor"
	HumiditySensor           = "HumiditySensor"
	LeakSensor               = "LeakSensor"
	Light                    = "Light"
	Lock                     = "Lock"
	MotionSensor             = "MotionSensor"
	MultiLevelSensor         = "MultiLevelSensor"
	MultiLevelSwitch         = "MultiLevelSwitch"
	OnOffSwitch              = "OnOffSwitch"
	SmartPlug                = "SmartPlug"
	SmokeSensor              = "SmokeSensor"
	TemperatureSensor        = "TemperatureSensor"
	Thermostat               = "Thermostat"
	VideoCamera              = "VideoCamera"

	Context = "https://webthings.io/schemas"

	OpenProperty              = "OpenProperty"
	PushedProperty            = "PushedProperty"
	SmokeProperty             = "SmokeProperty"
	TargetTemperatureProperty = "TargetTemperatureProperty"
	TemperatureProperty       = "TemperatureProperty"
	ThermostatModeProperty    = "ThermostatModeProperty"
	VideoProperty             = "VideoProperty"
	VoltageProperty           = "VoltageProperty"
)

const (
	AdapterAddedNotification           = 4096
	AdapterCancelPairingCommand        = 4100
	AdapterCancelRemoveDeviceCommand   = 4105
	AdapterPairingPromptNotification   = 4101
	AdapterRemoveDeviceRequest         = 4103
	AdapterRemoveDeviceResponse        = 4104
	AdapterStartPairingCommand         = 4099
	AdapterUnloadRequest               = 4097
	AdapterUnloadResponse              = 4098
	AdapterUnpairingPromptNotification = 4102
	ApiHandlerAddedNotification        = 20480
	ApiHandlerApiRequest               = 20483
	ApiHandlerApiResponse              = 20484
	ApiHandlerUnloadRequest            = 20481
	ApiHandlerUnloadResponse           = 20482
	DeviceActionStatusNotification     = 8201
	DeviceAddedNotification            = 8192
	DeviceConnectedStateNotification   = 8197
	DeviceDebugCommand                 = 8206
	DeviceEventNotification            = 8200
	DevicePropertyChangedNotification  = 8199
	DeviceRemoveActionRequest          = 8202
	DeviceRemoveActionResponse         = 8203
	DeviceRequestActionRequest         = 8204
	DeviceRequestActionResponse        = 8205
	DeviceSavedNotification            = 8207
	DeviceSetCredentialsRequest        = 8195
	DeviceSetCredentialsResponse       = 8196
	DeviceSetPinRequest                = 8193
	DeviceSetPinResponse               = 8194
	DeviceSetPropertyCommand           = 8198
	MockAdapterAddDeviceRequest        = 61440
	MockAdapterAddDeviceResponse       = 61441
	MockAdapterClearStateRequest       = 61446
	MockAdapterClearStateResponse      = 61447
	MockAdapterPairDeviceCommand       = 61444
	MockAdapterRemoveDeviceRequest     = 61442
	MockAdapterRemoveDeviceResponse    = 61443
	MockAdapterUnpairDeviceCommand     = 61445
	NotifierAddedNotification          = 12288
	NotifierUnloadRequest              = 12289
	NotifierUnloadResponse             = 12290
	OutletAddedNotification            = 16384
	OutletNotifyRequest                = 16386
	OutletNotifyResponse               = 16387
	OutletRemovedNotification          = 16385
	PluginErrorNotification            = 4
	PluginRegisterRequest              = 0
	PluginRegisterResponse             = 1
	PluginUnloadRequest                = 2
	PluginUnloadResponse               = 3
)

func MessageTypeToString(mt int) string {
	switch mt {

	case AdapterAddedNotification:
		return "AdapterAddedNotification"
	case AdapterCancelPairingCommand:
		return "AdapterCancelPairingCommand"
	case AdapterCancelRemoveDeviceCommand:
		return "AdapterCancelRemoveDeviceCommand"
	case AdapterPairingPromptNotification:
		return "AdapterPairingPromptNotification"
	case AdapterRemoveDeviceRequest:
		return "AdapterRemoveDeviceRequest"
	case AdapterRemoveDeviceResponse:
		return "AdapterRemoveDeviceResponse"
	case AdapterStartPairingCommand:
		return "AdapterStartPairingCommand"
	case AdapterUnloadRequest:
		return "AdapterUnloadRequest"
	case AdapterUnloadResponse:
		return "AdapterUnloadResponse"
	case AdapterUnpairingPromptNotification:
		return "AdapterUnpairingPromptNotification"
	case ApiHandlerAddedNotification:
		return "ApiHandlerAddedNotification"
	case ApiHandlerApiRequest:
		return "ApiHandlerApiRequest"
	case ApiHandlerApiResponse:
		return "ApiHandlerApiResponse"
	case ApiHandlerUnloadRequest:
		return "ApiHandlerUnloadRequest"
	case ApiHandlerUnloadResponse:
		return "ApiHandlerUnloadResponse"
	case DeviceActionStatusNotification:
		return "DeviceActionStatusNotification"
	case DeviceAddedNotification:
		return "DeviceAddedNotification"
	case DeviceConnectedStateNotification:
		return "DeviceConnectedStateNotification"
	case DeviceDebugCommand:
		return "DeviceDebugCommand"
	case DeviceEventNotification:
		return "DeviceEventNotification"
	case DevicePropertyChangedNotification:
		return "DevicePropertyChangedNotification"
	case DeviceRemoveActionRequest:
		return "DeviceRemoveActionRequest"
	case DeviceRemoveActionResponse:
		return "DeviceRemoveActionResponse"
	case DeviceRequestActionRequest:
		return "DeviceRequestActionRequest"
	case DeviceRequestActionResponse:
		return "DeviceRequestActionResponse"
	case DeviceSavedNotification:
		return "DeviceSavedNotification"
	case DeviceSetCredentialsRequest:
		return "DeviceSetCredentialsRequest"
	case DeviceSetCredentialsResponse:
		return "DeviceSetCredentialsResponse"
	case DeviceSetPinRequest:
		return "DeviceSetPinRequest"
	case DeviceSetPinResponse:
		return "DeviceSetPinResponse"
	case DeviceSetPropertyCommand:
		return "DeviceSetPropertyCommand"
	case MockAdapterAddDeviceRequest:
		return "MockAdapterAddDeviceRequest"
	case MockAdapterAddDeviceResponse:
		return "MockAdapterAddDeviceResponse"
	case MockAdapterClearStateRequest:
		return "MockAdapterClearStateRequest"
	case MockAdapterClearStateResponse:
		return "MockAdapterClearStateResponse"
	case MockAdapterPairDeviceCommand:
		return "MockAdapterPairDeviceCommand"
	case MockAdapterRemoveDeviceRequest:
		return "MockAdapterRemoveDeviceRequest"
	case MockAdapterRemoveDeviceResponse:
		return "MockAdapterRemoveDeviceResponse"
	case MockAdapterUnpairDeviceCommand:
		return "MockAdapterUnpairDeviceCommand"
	case NotifierAddedNotification:
		return "NotifierAddedNotification"
	case NotifierUnloadRequest:
		return "NotifierUnloadRequest"
	case NotifierUnloadResponse:
		return "NotifierUnloadResponse"
	case OutletAddedNotification:
		return "OutletAddedNotification"
	case OutletNotifyRequest:
		return "OutletNotifyRequest"
	case OutletNotifyResponse:
		return "OutletNotifyResponse"
	case OutletRemovedNotification:
		return "OutletRemovedNotification"
	case PluginErrorNotification:
		return "PluginErrorNotification"
	case PluginRegisterRequest:
		return "PluginRegisterRequest"
	case PluginRegisterResponse:
		return "PluginRegisterResponse"
	case PluginUnloadRequest:
		return "PluginUnloadRequest"
	case PluginUnloadResponse:
		return "PluginUnloadResponse"
	default:
		return "unknown"
	}
}

type Map = map[string]interface{}

const IpcDefaultPort = "9500"
