package config

import (
	_ "embed"
	"fmt"
	"gateway/log"
	"gateway/pkg/database"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	UnitCelsius = "celsius"

	preferencesKey = "settings.preferences"
)

var Conf *Config
var preferences *Preferences
var userProfile *UserProfile

//go:embed default.yaml
var DefaultConfig []byte

type addonManger struct {
	ListUrls []string `yaml:"listUrls"`
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
	Ports map[string]int `yaml:"ports"`

	AddonManager addonManger `yaml:"addonManager"`

	RemoveBeforeOpen bool   `yaml:"removeBeforeOpen"`
	LogRotateDays    int    `yaml:"logRotateDays"`
	NodeLoader       string `yaml:"nodeLoader"`

	ProfileDir string `yaml:"profileDir"`
	AddonsDir  string `yaml:"addonDir,omitempty"`
	LogDir     string `yaml:"logDir,omitempty"`
	DataDir    string `yaml:"dataDir,omitempty"`
	ConfigDir  string `yaml:"configDir,omitempty"`
	MediaDir   string `yaml:"mediaDir,omitempty"`
	UploadDir  string `yaml:"uploadDir,omitempty"`

	Architecture   string
	GatewayVersion string `yaml:"gatewayVersion,omitempty"`
}

func GetAddonListUrls() []string {
	return Conf.AddonManager.ListUrls
}

func InitRuntime(config string) error {

	var data []byte
	var rtc Config
	var err error
	if config != "" {
		data, err = ioutil.ReadFile(config)
	} else {
		data = DefaultConfig
	}
	err = yaml.Unmarshal(data, &rtc)
	if err != nil {
		return err
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
		return err
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
	log.InitLogger(rtc.LogDir, true, rtc.LogRotateDays)

	log.Info(fmt.Sprintf("gateway start path: %s", rtc.ProfileDir))

	//init database
	if rtc.RemoveBeforeOpen {
		database.ResetDB(rtc.ConfigDir)
	}
	err = database.InitDB(rtc.ConfigDir)

	if err != nil {
		return err
	}
	Conf = &rtc
	err = UpdateOrCreatePreferences()
	return err
}

func GetUserProfile() *UserProfile {
	return userProfile
}

func GetPreferences() *Preferences {
	return preferences
}
