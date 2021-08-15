package internal

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


