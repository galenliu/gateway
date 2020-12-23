package config

import (
	_ "embed"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

var Conf *Config

//go:embed default.yaml
var DefaultConfig []byte

type addonManger struct {
	ListUrls []string `yaml:"listUrls"`
}

var preferences *Preferences
var userProfile *UserProfile

type UserProfile struct {
	BaseDir        string `validate:"required" json:"base_dir"`
	DataDir        string `validate:"required" json:"data_dir"`
	AddonsDir      string `validate:"required" json:"addons_dir"`
	ConfigDir      string `validate:"required" json:"config_dir"`
	UploadDir      string `validate:"required" json:"upload_dir"`
	MediaDir       string `validate:"required" json:"media_dir"`
	LogDir         string `validate:"required" json:"log_dir"`
	GatewayVersion string `json:"gateway_version"`
}

type Units struct {
	Temperature string `gorm:"default: degree_celsius" json:"temperature"`
}

type Preferences struct {
	Language string `gorm:"default: zh-cn" json:"language"`
	Units    Units  `json:"units"`
	UnitsID  int    `json:"-"`
}

func GetUserProfile() *UserProfile {
	return userProfile
}

func GetPreferences() *Preferences {
	return preferences
}

func UpdatePreferences() *Preferences {
	//open database and create table
	db := database.GetDB()
	_ = db.AutoMigrate(&Preferences{})
	_ = db.AutoMigrate(&Units{})

	var pref Preferences
	result := db.First(&pref)
	if result.Error != nil {
		u1 := Units{Temperature: util.PrefUnitsTempCelsius}
		pref = Preferences{Language: util.PrefLangCn}
		pref.Units = u1
		db.Debug().Create(&pref)
		_ = db.First(&pref)
	}
	return &pref
}

type Config struct {
	Ports map[string]int `yaml:"ports"`

	AddonManager addonManger `yaml:"addonManager"`

	RemoveBeforeOpen bool `yaml:"removeBeforeOpen"`
	LogRotateDays    int  `yaml:"logRotateDays"`

	ProfileDir string `yaml:"profileDir"`
	AddonsDir  string `yaml:"addonDir,omitempty"`
	LogDir     string `yaml:"logDir,omitempty"`
	DataDir    string `yaml:"dataDir,omitempty"`
	ConfigDir  string `yaml:"configDir,omitempty"`
	MediaDir   string `yaml:"mediaDir,omitempty"`
	UploadDir  string `yaml:"uploadDir,omitempty"`

	Architecture string

	GatewayVersion string `yaml:"_"`
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
	if rtc.ProfileDir == "" {
		rtc.ProfileDir = util.GetDefaultConfigDir()
	}
	if rtc.AddonsDir == "" {
		rtc.AddonsDir = path.Join(rtc.ProfileDir, util.AddonsDir)
	}
	if rtc.LogDir == "" {
		rtc.LogDir = path.Join(rtc.ProfileDir, util.LogDir)
	}
	if rtc.DataDir == "" {
		rtc.DataDir = path.Join(rtc.ProfileDir, util.DataDir)
	}
	if rtc.ConfigDir == "" {
		rtc.ConfigDir = path.Join(rtc.ProfileDir, util.ConfigDir)
	}
	if rtc.MediaDir == "" {
		rtc.MediaDir = path.Join(rtc.ProfileDir, util.MediaDir)
	}
	if rtc.UploadDir == "" {
		rtc.UploadDir = path.Join(rtc.ProfileDir, util.UploadDir)
	}
	if rtc.GatewayVersion == "" {
		rtc.GatewayVersion = util.Version
	}
	err = util.EnsureDir(rtc.ProfileDir, rtc.AddonsDir, rtc.LogDir, rtc.DataDir, rtc.ConfigDir, rtc.MediaDir, rtc.UploadDir)
	if err != nil {
		return err
	}

	userProfile = &UserProfile{
		BaseDir:        rtc.ProfileDir,
		DataDir:        rtc.DataDir,
		AddonsDir:      rtc.AddonsDir,
		ConfigDir:      rtc.ConfigDir,
		UploadDir:      rtc.UploadDir,
		MediaDir:       rtc.MediaDir,
		LogDir:         rtc.LogDir,
		GatewayVersion: rtc.GatewayVersion,
	}
	//init logger
	log.InitLogger(rtc.LogDir, true, rtc.LogRotateDays)

	//init database
	if rtc.RemoveBeforeOpen {
		database.ResetDB(rtc.ConfigDir)
	}
	err = database.InitDB(rtc.ConfigDir)
	if err != nil {
		return err
	}
	Conf = &rtc

	preferences = UpdatePreferences()

	return err
}
