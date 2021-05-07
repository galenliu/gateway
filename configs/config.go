package configs

import (
	_ "embed"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"sync"

	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	UnitCelsius    = "celsius"
	preferencesKey = "settings.preferences"
)

var instance *Config
var preferences *Preferences
var userProfile *UserProfile
var once sync.Once

//go:embed default.json
var DefaultConfig []byte

type addonManger struct {
	ListUrls []string `json:"listUrls"`
}

type UserProfile struct {
	BaseDir    string `json:"baseDir"`
	DataDir    string `json:"dataDir"`
	AddonsDir  string `json:"addonsDir"`
	ConfigDir  string `json:"configDir"`
	UploadDir  string `json:"uploadDir"`
	MediaDir   string `json:"mediaDir"`
	LogDir     string `json:"logDir"`
	GatewayDir string `json:"gatewayDir"`
}

type Preferences struct {
	Language string `yaml:"language" json:"language"`
	Units    struct {
		Temperature string `json:"temperature"`
	} `json:"units"`
}

func preferencesDefault() *Preferences {
	return &Preferences{
		Language: "zh-cn",
		Units: struct {
			Temperature string `json:"temperature"`
		}{Temperature: UnitCelsius},
	}
}

func UpdateOrCreatePreferences() error {
	p, err := database.GetSetting(preferencesKey)
	if err != nil {
		s, ee := json.MarshalToString(preferencesDefault())
		if ee != nil {
			return ee
		}
		eee := database.SetSetting(preferencesKey, s)
		if eee != nil {
			return eee
		}
		p, _ = database.GetSetting(preferencesKey)
	}
	_ = json.UnmarshalFromString(p, &preferences)
	return nil
}

type Config struct {
	Ports          Ports        `json:"ports"`
	ProfileDir     string       `json:"profile_dir"`
	GatewayVersion string       `json:"gateway_version"`
	AddonManager   AddonManager `json:"addon_manager"`
	Log            Log          `json:"log"`
	Database       Database     `json:"database"`

	AddonsDir    string `json:"addonDir,omitempty"`
	LogDir       string `json:"logDir,omitempty"`
	DataDir      string `json:"dataDir,omitempty"`
	ConfigDir    string `json:"configDir,omitempty"`
	MediaDir     string `json:"mediaDir,omitempty"`
	UploadDir    string `json:"uploadDir,omitempty"`
	Architecture string
}

type Ports struct {
	HTTPS int `json:"https"`
	HTTP  int `json:"http"`
	Ipc   int `json:"ipc"`
}
type AddonManager struct {
	NodeLoader string   `json:"node_loader"`
	ListUrls   []string `json:"list_urls"`
	TestAddons bool     `json:"test_addons"`
}
type Log struct {
	Verbose       bool `json:"verbose"`
	LogRotateDays int  `json:"log_rotate_days"`
}
type Database struct {
	RemoveBeforeOpen bool `json:"remove_before_open"`
}

func GetAddonListUrls() []string {
	return instance.AddonManager.ListUrls
}

// NewConfig 加载配置文件，初始化Log,DataBase,工作目录
func NewConfig(config string) *Config {
	once.Do(func() {
		var data []byte
		var rtc Config
		var err error
		if config != "" {
			data, err = ioutil.ReadFile(config)
		} else {
			data = DefaultConfig
		}
		err = json.Unmarshal(data, &rtc)
		if err != nil {
			return
		}
		if !path.IsAbs(rtc.ProfileDir) {
			rtc.ProfileDir, _ = filepath.Abs(rtc.ProfileDir)
		}

		if rtc.ProfileDir == "" {
			rtc.ProfileDir = util.GetDefaultConfigDir()
		}
		if rtc.AddonsDir == "" {
			rtc.AddonsDir = rtc.ProfileDir + string(os.PathSeparator) + util.AddonsDir
		}
		if rtc.LogDir == "" {
			rtc.LogDir = rtc.ProfileDir + string(os.PathSeparator) + util.LogDir

		}
		if rtc.DataDir == "" {
			rtc.DataDir = rtc.ProfileDir + string(os.PathSeparator) + util.DataDir
		}
		if rtc.ConfigDir == "" {
			rtc.ConfigDir = rtc.ProfileDir + string(os.PathSeparator) + util.ConfigDir
		}
		if rtc.MediaDir == "" {
			rtc.MediaDir = rtc.ProfileDir + string(os.PathSeparator) + util.MediaDir
		}
		if rtc.UploadDir == "" {
			rtc.UploadDir = rtc.ProfileDir + string(os.PathSeparator) + util.UploadDir
		}
		if rtc.GatewayVersion == "" {
			rtc.GatewayVersion = util.Version
		}
		err = util.EnsureDir(rtc.ProfileDir, rtc.AddonsDir, rtc.LogDir, rtc.DataDir, rtc.ConfigDir, rtc.MediaDir, rtc.UploadDir)
		if err != nil {
			return
		}

		userProfile = &UserProfile{
			BaseDir:    rtc.ProfileDir,
			DataDir:    rtc.DataDir,
			AddonsDir:  rtc.AddonsDir,
			ConfigDir:  rtc.ConfigDir,
			UploadDir:  rtc.UploadDir,
			MediaDir:   rtc.MediaDir,
			LogDir:     rtc.LogDir,
			GatewayDir: rtc.ProfileDir,
		}

		//init logger
		log.InitLogger(rtc.LogDir, true, rtc.Log.LogRotateDays)

		//init database
		if rtc.Database.RemoveBeforeOpen {
			database.ResetDB(rtc.ConfigDir)
		}
		err = database.InitDB(rtc.ConfigDir)

		log.Info("gateway loaded config on path: %s", rtc.ProfileDir)

		if err != nil {
			return
		}
		instance = &rtc
		err = UpdateOrCreatePreferences()
		return
	})

	return instance
}

func GetUserProfile() *UserProfile {
	return userProfile
}

func GetPreferences() *Preferences {
	return preferences
}

func IsVerbose() bool {
	return instance.Log.Verbose
}

func GetIpcPort() int {
	return instance.Ports.Ipc
}

func GetAddonsDir() string {
	return instance.AddonsDir
}

func GetDataDir() string {
	return instance.DataDir
}

func GetLogDir() string {
	return instance.LogDir
}

func GetArchitecture() string {
	return instance.Architecture
}

func GetConfigDir() string {
	return instance.ConfigDir
}

func GetGatewayVersion() string {
	return instance.GatewayVersion
}

func GetProfileDir() string {
	return instance.ProfileDir
}

func GetUploadDir() string {
	return instance.UploadDir
}
func GetNodeLoader() string {
	return instance.AddonManager.NodeLoader
}

func GetPorts() Ports {
	return instance.Ports
}
