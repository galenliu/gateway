package config

import (
	_ "embed"
	"errors"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"io/ioutil"
	"path"
	"time"
)

const (
	UnitCelsius = "celsius"
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
	gorm.Model
	Temperature   string
	PreferencesID string
}

type Preferences struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"primarykey"`
	Language  string
	Units     Units `gorm:"foreignKey:PreferencesID"`
}

func UpdateOrCreatePreferences(pref *Preferences) error {

	db, err := database.GetDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Preferences{}, &Units{})
	if err != nil {
		return err
	}
	if pref != nil {
		tx := db.Save(&pref)
		return tx.Error
	}
	var p Preferences
	rst := db.Where("name = ?", "preferences").First(&p)
	if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
		p = Preferences{
			Name: "preferences",
			Language: "zh-cn",
			Units:    Units{Temperature: UnitCelsius},
		}
		db.Create(&p)
	}
	preferences = &p
	return nil
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

	Architecture   string
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
	err = UpdateOrCreatePreferences(nil)
	return err
}

func GetUserProfile() *UserProfile {
	return userProfile
}

func GetPreferences() *Preferences {
	return preferences
}
