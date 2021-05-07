package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

const (
	// Web server routes
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

	// Plugin and REST/websocket API things
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
	REMOVED         = "removed"
	RequestAction   = "requestAction"
	SetProperty     = "setProperty"
	//ThingAdded event args: []byte
	ThingAdded    = "thingAdded"
	ThingModified = "thingModified"
	ThingRemoved  = "thingRemoved"

	WebServerStarted    = "webServerStarted"
	PluginServerStarted = "pluginServerStarted"
	WebServerStopped    = "webServerStopped"
	PluginServerStopped = "pluginServerStopped"

	// OAuth things
	AccessToken       = "access_token"
	AuthorizationCode = "authorization_code"
	UserToken         = "user_token"
	READWRITE         = "readwrite"
	READ              = "read"

	MajorVersion = 1
	MinorVersion = 0
	PatchVersion = 0

	AddonsDir = "addons"
	DataDir   = "data"
	ConfigDir = "config"
	MediaDir  = "media"
	UploadDir = "upload"
	LogDir    = "logger"

	DbPrefLang = "preferences.language"
	PrefLangCn = "zh-CN"

	DbPrefUnitsTemp      = "preferences.units.temperature"
	PrefUnitsTempCelsius = "degree celsius"
	ConfDirName          = ".gateway"
)

var (
	ShortVersion = fmt.Sprintf("%v.%v", MajorVersion, MinorVersion)
	Version      = fmt.Sprintf("%v.%v", ShortVersion, PatchVersion)
)

type Map = fiber.Map
