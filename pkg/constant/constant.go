package constant

// events
const (
	Support = ""

	UnloadPluginKillDelay = 3000

	UsersPath        = "/users"
	ThingsPath       = "/things"
	PropertiesPath   = "/properties"
	NewThingsPath    = "/new_things"
	AdaptersPath     = "/adapters"
	AddonsPath       = "/addons"
	NotifiersPath    = "/notifiers"
	ActionsPath      = "actions"
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
	GroupsPath       = "/groups"
	OauthclientsPath = "/authorizations"
	InternalLogsPath = "/internal-logs"
	LogsPath         = "/logs"
	PushPath         = "/push"
	PingPath         = "/ping"
	ProxyPath        = "/proxy"
	ExtensionsPath   = "/extensions"
)

const (
	ActionStatus         = "actionStatus"
	AdapterAdded         = "adapterAdded"
	AddEventSubscription = "addEventSubscription"
	ApiHandlerAdded      = "apiHandlerAdded"

	Connected = "connected"
	Removed   = "Removed"
	Created   = "created"
	Event     = "events"
	Modified  = "modified"

	ThingModified = "thingModified"
	ThingRemoved  = "thingRemoved"
	ThingAdded    = "thingAdded"
	ThingCreated  = "thingCreated"

	Completed       = "completed"
	ERROR           = "error"
	NotifierAdded   = "notifierAdded"
	OutletAdded     = "outletAdded"
	OutletRemoved   = "outletRemoved"
	PairingTimeout  = "pairingTimeout"
	PropertyChanged = "propertyChanged"
	PropertyStatus  = "propertyStatus"
	RequestAction   = "requestAction"
	SetProperty     = "setProperty"

	// WebServerStarted Web api events
	WebServerStarted  = "webServerStarted"
	WebServerStopped  = "webServerStopped"
	AccessToken       = "access_token"
	AuthorizationCode = "authorization_code"
	UserToken         = "user_token"
	READWRITE         = "readwrite"
	READ              = "read"

	PrefUnitsTempCelsius = "degree celsius"
	DegreeCelsius        = "degree celsius"    //摄氏度
	FahrenheitDegree     = "fahrenheit degree" //华氏度

	DeviceRemovalTimeout = 30000
)

const (
	ZhCN = "zh-CN"  //简体中文(中国)
	ZhHk = "zh-HK"  //繁体中文(香港)
	ZhTW = "zh-TW"  //繁体中文(台湾地区)
	EnHk = "en-HK"  //英语(香港)
	EnUS = " en-US" //英语(美国)
	EnGB = "en-GB"  //英语(英国)
	EnWW = "en-WW"  //英语(全球)
	EnCA = "en-CA"  //英语(加拿大)
	EnAU = "en-AU"  //英语(澳大利亚)
	EnIE = "en-IE"  //英语(爱尔兰)
	EnFI = "en-FI"  //英语(芬兰)
	FiFI = "fi-FI"  //芬兰语(芬兰)
	EnDK = "en-DK"  //英语(丹麦)
	DaDK = "da-DK"  //丹麦语(丹麦)
	EnIL = "en-IL"  //英语(以色列)
	HeIL = "he-IL"  //希伯来语(以色列)
	EnZA = "en-ZA"  //英语(南非)
	EnIN = "en-IN"  //英语(印度)
	EnNO = "en-NO"  //英语(挪威)
	EnSG = "en-SG"  //英语(新加坡)
	EnNZ = "en-NZ"  //英语(新西兰)
	EnID = "en-id"  //英语(印度尼西亚)
	EnPH = "en-PH"  //英语(菲律宾)
	EnTH = "en-TH"  //英语(泰国)
	EnMY = "en-MY"  // 英语(马来西亚)
	EnXA = "en-XA"  //英语(阿拉伯)
	KoKR = "ko-KR"  //韩文(韩国)
	JaJP = "ja-JP"  //日语(日本)
	NlNL = "nl-NL"  //荷兰语(荷兰)
	NlBE = "nl-BE"  //荷兰语(比利时)
	PtPT = "pt-PT"  //葡萄牙语(葡萄牙)
	PtBR = "pt-BR"  //葡萄牙语(巴西)
	FrFR = "fr-FR " //法语(法国)
	FrLU = "fr-LU"  //法语(卢森堡)
	FrCH = "fr-CH"  //法语(瑞士)
	FrBE = "fr-BE"  //法语(比利时)
	FrCA = "fr-CA"  //法语(加拿大)
	EsLA = "es-LA"  //西班牙语(拉丁美洲)
	EsES = "es-ES"  //西班牙语(西班牙)
	EsAR = "es-AR"  //西班牙语(阿根廷)
	EsUS = "es-US"  //西班牙语(美国)
	EsMX = "es-MX"  //西班牙语(墨西哥)
	EsCO = "es-CO"  //西班牙语(哥伦比亚)
	EsPR = "es-PR"  //西班牙语(波多黎各)
	DeDE = " de-DE" //德语(德国)
	DeAT = "de-AT"  //德语(奥地利)
	DeCH = "de-CH"  //德语(瑞士)
	RuRU = "ru-RU"  //俄语(俄罗斯)
	ItIT = "it-IT"  //意大利语(意大利)
	ElGR = "el-GR"  //希腊语(希腊)
	NoNO = "no-NO"  //挪威语(挪威)
	HuHU = "hu-HU"  //匈牙利语(匈牙利)
	TrTR = "tr-TR"  //土耳其语(土耳其)
	CsCZ = "cs-CZ"  //捷克语(捷克共和国)
	SlSL = "sl-SL"  //斯洛文尼亚语
	PlPL = "pl-PL"  //波兰语(波兰)
	SvSE = "sv-SE"  //瑞典语(瑞典)
)
