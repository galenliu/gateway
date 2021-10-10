package constant

// Events
const (
	GatewayStart = "gatewayStart"
	GatewayStop  = "gatewayStop"

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
	ServicesPath     = "/services"
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
	ServiceAdded         = "serviceAdded"
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

	DbPrefUnitsTemp      = "preferences.units.temperature"
	PrefUnitsTempCelsius = "degree celsius"
	ConfDirName          = ".gateway"
)
