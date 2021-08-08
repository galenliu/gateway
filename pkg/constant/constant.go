package constant

import (
	"github.com/gofiber/fiber/v2"
)

// Events
const (
	GatewayStarted = "gatewayStarted"
	GatewayStopped = "gatewayStopped"

	ThingAdded    = "thingAdded"
	ThingModified = "thingModified"
	ThingRemoved  = "thingRemoved"
	ThingCreated  = "thingCreated"

	// DeviceAdded Addon Manager Event
	DeviceAdded         = "deviceAdded"
	DeviceRemoved       = "deviceRemoved"
	AddonManagerStarted = "addonManagerStarted"
	AddonManagerStopped = "addonManagerStopped"
)

// Web server routes
const (
	UsersPath        = "/users"
	ThingsPath       = "/things"
	PropertiesPath   = "/properties"
	NewThingsPath    = "/new_things"
	AdaptersPath     = "/adapters"
	AddonsPath       = "/addons"
	NotifiersPath    = "/notifiers"
	ActionsPath      = "/actions"
	EventsPath       = "/events"
	LoginPath        = "/login"
	LogOutPath       = "/log-out"
	SettingsPath     = "/settings"
	UpdatesPath      = "/updates"
	UploadsPath      = "/uploads"
	MediaPath        = "/media"
	DebugPath        = "/debug"
	RulesPath        = "/rules"
	OauthPath        = "/oauth"
	OauthclientsPath = "/authorizations"
	InternalLogsPath = "/internal-logs"
	LogsPath         = "/logs"
	PushPath         = "/push"
	PingPath         = "/ping"
	ProxyPath        = "/proxy"
	ExtensionsPath   = "/extensions"
	ThingIdParam     = "/:thingId"

	ActionStatus         = "actionStatus"
	AdapterAdded         = "adapterAdded"
	AddEventSubscription = "addEventSubscription"
	ApiHandlerAdded      = "apiHandlerAdded"
	CONNECTED            = "connected"

	Created         = "created"
	Completed       = "completed"
	ERROR           = "error"
	EVENT           = "event"
	MODIFIED        = "modified"
	NotifierAdded   = "notifierAdded"
	OutletAdded     = "outletAdded"
	OutletRemoved   = "outletRemoved"
	PairingTimeout  = "pairingTimeout"
	PropertyChanged = "propertyChanged"
	PropertyStatus  = "propertyStatus"
	RequestAction   = "requestAction"
	SetProperty     = "setProperty"

	// WebServerStarted Web server event
	WebServerStarted  = "webServerStarted"
	WebServerStopped  = "webServerStopped"
	AccessToken       = "access_token"
	AuthorizationCode = "authorization_code"
	UserToken         = "user_token"
	READWRITE         = "readwrite"
	READ              = "read"

	AddonsDirName = "addons"
	DataDirName   = "data"
	ConfigDirName = "config"
	MediaDirName  = "media"
	UploadDirName = "upload"
	LogDirName    = "logger"

	DbPrefLang = "preferences.language"
	PrefLangCn = "zh-CN"

	DbPrefUnitsTemp      = "preferences.units.temperature"
	PrefUnitsTempCelsius = "degree celsius"
	ConfDirName          = ".gateway"
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

type Map = fiber.Map
