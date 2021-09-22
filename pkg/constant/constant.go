package constant

// events
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

	DbPrefLang = "preferences.language"
	PrefLangCn = "zh-CN"

	DbPrefUnitsTemp      = "preferences.units.temperature"
	PrefUnitsTempCelsius = "degree celsius"
	ConfDirName          = ".gateway"

	DegreeCelsius    = "degree celsius"    //摄氏度
	FahrenheitDegree = "fahrenheit degree" //华氏度
)

const (
	ZhCN = "zh-cn"  //简体中文(中国)
	ZhHk = "zh-hk"  //繁体中文(香港)
	ZhTW = "zh-tw"  //繁体中文(台湾地区)
	EnHk = "en-hk"  //英语(香港)
	EnUS = " en-us" //英语(美国)
	EnGB = "en-gb"  //英语(英国)
	EnWW = "en-ww"  //英语(全球)
	EnCA = "en-ca"  //英语(加拿大)
	EnAU = "en-au"  //英语(澳大利亚)
	EnIE = "en-ie"  //英语(爱尔兰)
	EnFI = "en-fi"  //英语(芬兰)
	FiFI = "fi-fi"  //芬兰语(芬兰)
	EnDK = "en-dk"  //英语(丹麦)
	DaDK = "da-dk"  //丹麦语(丹麦)
	EnIL = "en-il"  //英语(以色列)
	HeIL = "he-il"  //希伯来语(以色列)
	EnZA = "en-za"  //英语(南非)
	EnIN = "en-in"  //英语(印度)
	EnNO = "en-no"  //英语(挪威)
	EnSG = "en-sg"  //英语(新加坡)
	EnNZ = "en-nz"  //英语(新西兰)
	EnID = "en-id"  //英语(印度尼西亚)
	EnPH = "en-ph"  //英语(菲律宾)
	EnTH = "en-th"  //英语(泰国)
	EnMY = "en-my"  // 英语(马来西亚)
	EnXA = "en-xa"  //英语(阿拉伯)
	KoKR = "ko-kr"  //韩文(韩国)
	JaJP = "ja-jp"  //日语(日本)
	NlNL = "nl-nl"  //荷兰语(荷兰)
	NlBE = "nl-be"  //荷兰语(比利时)
	PtPT = "pt-pt"  //葡萄牙语(葡萄牙)
	PtBR = "pt-br"  //葡萄牙语(巴西)
	FrFR = "fr-fr " //法语(法国)
	FrLU = "fr-lu"  //法语(卢森堡)
	FrCH = "fr-ch"  //法语(瑞士)
	FrBE = "fr-be"  //法语(比利时)
	FrCA = "fr-ca"  //法语(加拿大)
	EsLA = "es-la"  //西班牙语(拉丁美洲)
	EsES = "es-es"  //西班牙语(西班牙)
	EsAR = "es-ar"  //西班牙语(阿根廷)
	EsUS = "es-us"  //西班牙语(美国)
	EsMX = "es-mx"  //西班牙语(墨西哥)
	EsCO = "es-co"  //西班牙语(哥伦比亚)
	EsPR = "es-pr"  //西班牙语(波多黎各)
	DeDE = " de-de" //德语(德国)
	DeAT = "de-at"  //德语(奥地利)
	DeCH = "de-ch"  //德语(瑞士)
	RuRU = "ru-ru"  //俄语(俄罗斯)
	ItIT = "it-it"  //意大利语(意大利)
	ElGR = "el-gr"  //希腊语(希腊)
	NoNO = "no-no"  //挪威语(挪威)
	HuHU = "hu-hu"  //匈牙利语(匈牙利)
	TrTR = "tr-tr"  //土耳其语(土耳其)
	CsCZ = "cs-cz"  //捷克语(捷克共和国)
	SlSL = "sl-sl"  //斯洛文尼亚语
	PlPL = "pl-pl"  //波兰语(波兰)
	SvSE = "sv-se"  //瑞典语(瑞典)
)
