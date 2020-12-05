package runtime

import (
	"gateway/pkg/database"
	"gateway/pkg/logger"
	"gateway/pkg/util"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

var RuntimeConf *RuntimeConfig

type addonManger struct {
	ListUrls []string `yaml:"listUrls"`
}

type RuntimeConfig struct {
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

var box *packr.Box

func GetAddonListUrls() []string {
	return RuntimeConf.AddonManager.ListUrls
}

func InitRuntime(config string) error {
	var data []byte
	var rtc RuntimeConfig
	var err error
	if config != "" {
		data, err = ioutil.ReadFile(config)
	} else {
		box = packr.New("config", "../../config")
		data, err = box.Find("default.yaml")
	}
	if err != nil {
		return err
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

	//init logger
	logger.InitLogger(rtc.LogDir, true, rtc.LogRotateDays)

	//init database
	if rtc.RemoveBeforeOpen {
		database.ResetDB(rtc.ConfigDir)
	}
	err = database.InitDB(rtc.ConfigDir)
	if err != nil {
		return err
	}
	RuntimeConf = &rtc

	return err
}
