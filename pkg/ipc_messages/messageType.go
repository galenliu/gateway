package messages

type MessageType int32

const (
	MessageType_PluginRegisterRequest              MessageType = 0
	MessageType_PluginRegisterResponse             MessageType = 1
	MessageType_PluginUnloadRequest                MessageType = 2
	MessageType_PluginUnloadResponse               MessageType = 3
	MessageType_PluginErrorNotification            MessageType = 4
	MessageType_AdapterAddedNotification           MessageType = 4096
	MessageType_AdapterCancelPairingCommand        MessageType = 4100
	MessageType_AdapterPairingPromptNotification   MessageType = 4101
	MessageType_AdapterUnpairingPromptNotification MessageType = 4102
	MessageType_AdapterRemoveDeviceRequest         MessageType = 4103
	MessageType_AdapterRemoveDeviceResponse        MessageType = 4104
	MessageType_AdapterCancelRemoveDeviceCommand   MessageType = 4105
	MessageType_AdapterUnloadRequest               MessageType = 4097
	MessageType_AdapterUnloadResponse              MessageType = 4098
	MessageType_AdapterStartPairingCommand         MessageType = 4099
	MessageType_ApiHandlerAddedNotification        MessageType = 20480
	MessageType_ApiHandlerApiRequest               MessageType = 20483
	MessageType_ApiHandlerApiResponse              MessageType = 20484
	MessageType_ApiHandlerUnloadRequest            MessageType = 20481
	MessageType_ApiHandlerUnloadResponse           MessageType = 20482
	MessageType_DeviceActionStatusNotification     MessageType = 8201
	MessageType_DeviceAddedNotification            MessageType = 8192
	MessageType_DeviceConnectedStateNotification   MessageType = 8197
	MessageType_DeviceDebugCommand                 MessageType = 8206
	MessageType_DeviceEventNotification            MessageType = 8200
	MessageType_DevicePropertyChangedNotification  MessageType = 8199
	MessageType_DeviceRemoveActionRequest          MessageType = 8202
	MessageType_DeviceRemoveActionResponse         MessageType = 8203
	MessageType_DeviceRequestActionRequest         MessageType = 8204
	MessageType_DeviceRequestActionResponse        MessageType = 8205
	MessageType_DeviceSavedNotification            MessageType = 8207
	MessageType_DeviceSetCredentialsRequest        MessageType = 8195
	MessageType_DeviceSetCredentialsResponse       MessageType = 8196
	MessageType_DeviceSetPinRequest                MessageType = 8193
	MessageType_DeviceSetPinResponse               MessageType = 8194
	MessageType_DeviceSetPropertyCommand           MessageType = 8198
	MessageType_MockAdapterAddDeviceRequest        MessageType = 61440
	MessageType_MockAdapterAddDeviceResponse       MessageType = 61441
	MessageType_MockAdapterClearStateRequest       MessageType = 61446
	MessageType_MockAdapterClearStateResponse      MessageType = 61447
	MessageType_MockAdapterPairDeviceCommand       MessageType = 61444
	MessageType_MockAdapterRemoveDeviceRequest     MessageType = 61442
	MessageType_MockAdapterRemoveDeviceResponse    MessageType = 61443
	MessageType_MockAdapterUnpairDeviceCommand     MessageType = 61445
	MessageType_NotifierAddedNotification          MessageType = 12288
	MessageType_NotifierUnloadRequest              MessageType = 12289
	MessageType_NotifierUnloadResponse             MessageType = 12290
	MessageType_OutletAddedNotification            MessageType = 16384
	MessageType_OutletNotifyRequest                MessageType = 16386
	MessageType_OutletNotifyResponse               MessageType = 16387
	MessageType_OutletRemovedNotification          MessageType = 16385
	MessageType_ServiceAddedNotification           MessageType = 81000
	MessageType_ServiceSetPropertyValueRequest     MessageType = 81100
	MessageType_ServicePropertyChangedNotification MessageType = 81101
	MessageType_ServiceActionsStatusNotification   MessageType = 81109
	MessageType_ServiceGetThingsRequest            MessageType = 81001
	MessageType_ServiceGetThingsResponse           MessageType = 81002
	MessageType_ServiceGetThingRequest             MessageType = 81003
	MessageType_ServiceGetThingResponse            MessageType = 81004
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0:     "PluginRegisterRequest",
		1:     "PluginRegisterResponse",
		2:     "PluginUnloadRequest",
		3:     "PluginUnloadResponse",
		4:     "PluginErrorNotification",
		4096:  "AdapterAddedNotification",
		4100:  "AdapterCancelPairingCommand",
		4101:  "AdapterPairingPromptNotification",
		4102:  "AdapterUnpairingPromptNotification",
		4103:  "AdapterRemoveDeviceRequest",
		4104:  "AdapterRemoveDeviceResponse",
		4105:  "AdapterCancelRemoveDeviceCommand",
		4097:  "AdapterUnloadRequest",
		4098:  "AdapterUnloadResponse",
		4099:  "AdapterStartPairingCommand",
		20480: "ApiHandlerAddedNotification",
		20483: "ApiHandlerApiRequest",
		20484: "ApiHandlerApiResponse",
		20481: "ApiHandlerUnloadRequest",
		20482: "ApiHandlerUnloadResponse",
		8201:  "DeviceActionStatusNotification",
		8192:  "DeviceAddedNotification",
		8197:  "DeviceConnectedStateNotification",
		8206:  "DeviceDebugCommand",
		8200:  "DeviceEventNotification",
		8199:  "DevicePropertyChangedNotification",
		8202:  "DeviceRemoveActionRequest",
		8203:  "DeviceRemoveActionResponse",
		8204:  "DeviceRequestActionRequest",
		8205:  "DeviceRequestActionResponse",
		8207:  "DeviceSavedNotification",
		8195:  "DeviceSetCredentialsRequest",
		8196:  "DeviceSetCredentialsResponse",
		8193:  "DeviceSetPinRequest",
		8194:  "DeviceSetPinResponse",
		8198:  "DeviceSetPropertyCommand",
		61440: "MockAdapterAddDeviceRequest",
		61441: "MockAdapterAddDeviceResponse",
		61446: "MockAdapterClearStateRequest",
		61447: "MockAdapterClearStateResponse",
		61444: "MockAdapterPairDeviceCommand",
		61442: "MockAdapterRemoveDeviceRequest",
		61443: "MockAdapterRemoveDeviceResponse",
		61445: "MockAdapterUnpairDeviceCommand",
		12288: "NotifierAddedNotification",
		12289: "NotifierUnloadRequest",
		12290: "NotifierUnloadResponse",
		16384: "OutletAddedNotification",
		16386: "OutletNotifyRequest",
		16387: "OutletNotifyResponse",
		16385: "OutletRemovedNotification",
		81000: "ServiceAddedNotification",
		81100: "ServiceSetPropertyValueRequest",
		81101: "ServicePropertyChangedNotification",
		81109: "ServiceActionsStatusNotification",
		81001: "ServiceGetThingsRequest",
		81002: "ServiceGetThingsResponse",
		81003: "ServiceGetThingRequest",
		81004: "ServiceGetThingResponse",
	}
	MessageType_value = map[string]int32{
		"PluginRegisterRequest":              0,
		"PluginRegisterResponse":             1,
		"PluginUnloadRequest":                2,
		"PluginUnloadResponse":               3,
		"PluginErrorNotification":            4,
		"AdapterAddedNotification":           4096,
		"AdapterCancelPairingCommand":        4100,
		"AdapterPairingPromptNotification":   4101,
		"AdapterUnpairingPromptNotification": 4102,
		"AdapterRemoveDeviceRequest":         4103,
		"AdapterRemoveDeviceResponse":        4104,
		"AdapterCancelRemoveDeviceCommand":   4105,
		"AdapterUnloadRequest":               4097,
		"AdapterUnloadResponse":              4098,
		"AdapterStartPairingCommand":         4099,
		"ApiHandlerAddedNotification":        20480,
		"ApiHandlerApiRequest":               20483,
		"ApiHandlerApiResponse":              20484,
		"ApiHandlerUnloadRequest":            20481,
		"ApiHandlerUnloadResponse":           20482,
		"DeviceActionStatusNotification":     8201,
		"DeviceAddedNotification":            8192,
		"DeviceConnectedStateNotification":   8197,
		"DeviceDebugCommand":                 8206,
		"DeviceEventNotification":            8200,
		"DevicePropertyChangedNotification":  8199,
		"DeviceRemoveActionRequest":          8202,
		"DeviceRemoveActionResponse":         8203,
		"DeviceRequestActionRequest":         8204,
		"DeviceRequestActionResponse":        8205,
		"DeviceSavedNotification":            8207,
		"DeviceSetCredentialsRequest":        8195,
		"DeviceSetCredentialsResponse":       8196,
		"DeviceSetPinRequest":                8193,
		"DeviceSetPinResponse":               8194,
		"DeviceSetPropertyCommand":           8198,
		"MockAdapterAddDeviceRequest":        61440,
		"MockAdapterAddDeviceResponse":       61441,
		"MockAdapterClearStateRequest":       61446,
		"MockAdapterClearStateResponse":      61447,
		"MockAdapterPairDeviceCommand":       61444,
		"MockAdapterRemoveDeviceRequest":     61442,
		"MockAdapterRemoveDeviceResponse":    61443,
		"MockAdapterUnpairDeviceCommand":     61445,
		"NotifierAddedNotification":          12288,
		"NotifierUnloadRequest":              12289,
		"NotifierUnloadResponse":             12290,
		"OutletAddedNotification":            16384,
		"OutletNotifyRequest":                16386,
		"OutletNotifyResponse":               16387,
		"OutletRemovedNotification":          16385,
		"ServiceAddedNotification":           81000,
		"ServiceSetPropertyValueRequest":     81100,
		"ServicePropertyChangedNotification": 81101,
		"ServiceActionsStatusNotification":   81109,
		"ServiceGetThingsRequest":            81001,
		"ServiceGetThingsResponse":           81002,
		"ServiceGetThingRequest":             81003,
		"ServiceGetThingResponse":            81004,
	}
)

type Type int32

const (
	Type_null    Type = 0
	Type_boolean Type = 1
	Type_object  Type = 2
	Type_array   Type = 3
	Type_number  Type = 4
	Type_integer Type = 5
	Type_string  Type = 6
)
