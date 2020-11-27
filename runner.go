package gateway

import (
	"gateway/addons"
	"gateway/app"
	"gateway/util"
	"gateway/util/database"
	"gateway/util/logger"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"runtime"
)

type addonManger struct {
	ListUrls   []string `yaml:"listUrls"`
	testAddons bool     `yaml:"testAddons"`
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

func InitRuntime(config string) (*RuntimeConfig, error) {
	var data []byte
	var rtc RuntimeConfig
	var err error
	if config != "" {
		data, err = ioutil.ReadFile(config)
	} else {
		box = packr.New("config", "./config")
		data, err = box.Find("default.yaml")
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &rtc)
	if err != nil {
		return nil, err
	}
	if rtc.ProfileDir == "" {
		rtc.ProfileDir = GetDefaultConfigDir()
	}
	if rtc.AddonsDir == "" {
		rtc.AddonsDir = path.Join(rtc.ProfileDir, AddonsDir)
	}
	if rtc.LogDir == "" {
		rtc.LogDir = path.Join(rtc.ProfileDir, LogDir)
	}
	if rtc.DataDir == "" {
		rtc.DataDir = path.Join(rtc.ProfileDir, DataDir)
	}
	if rtc.ConfigDir == "" {
		rtc.ConfigDir = path.Join(rtc.ProfileDir, ConfigDir)
	}
	if rtc.MediaDir == "" {
		rtc.MediaDir = path.Join(rtc.ProfileDir, MediaDir)
	}
	if rtc.UploadDir == "" {
		rtc.UploadDir = path.Join(rtc.ProfileDir, UploadDir)
	}
	if rtc.GatewayVersion == "" {
		rtc.GatewayVersion = Version
	}
	err = util.EnsureDir(rtc.ProfileDir, rtc.AddonsDir, rtc.LogDir, rtc.DataDir, rtc.ConfigDir, rtc.MediaDir, rtc.UploadDir)
	if err != nil {
		return nil, err
	}
	//init logger
	logger.InitLogger(rtc.LogDir, true, rtc.LogRotateDays)

	//init database
	if rtc.RemoveBeforeOpen {
		database.ResetDB(rtc.ConfigDir)
	}
	err = database.InitDB(rtc.ConfigDir)
	if err != nil {
		log.Error("database init err")
		return nil, err
	}

	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	rtc.Architecture = runtime.GOOS + "-" + arch

	return &rtc, err
}

func InitGateway(rtc *RuntimeConfig) (gateway *HomeGateway, err error) {

	gateway = &HomeGateway{}
	gateway.Runtime = rtc
	//update the gateway preferences
	gateway.updatePreferences()
	return gateway, err
}

func StartAddonsManager(gateway *HomeGateway) {
	config := addons.ManagerConfig{
		AddonsDir: gateway.Runtime.AddonsDir,
		DataDir:   gateway.Runtime.DataDir,
	}
	if gateway.AddonsManager == nil {
		gateway.AddonsManager = addons.NewAddonsManager(config)
	}
	if gateway.AddonsManager.IsRunning {
		gateway.AddonsManager.Stop()
	}

	gateway.AddonsManager.LoadAddons()

	//_ = gateway.AddonsManager.InstallAddonFromUrl("tuya-adapter",
	//	"https://gitee.com/liu_guilin/tuya-adapter/attach_files/525074/download/tuya-adapter-0.2.4.tgz",
	//	"2905594a1893443385c4f1cd5ed254bbdd4022b5e87520212e5a7cd8c9d0ab25", true)
	//_ = gateway.AddonsManager.InstallAddonFromUrl("tplink-adapter",
	//	"https://gitee.com/liu_guilin/tplink-adapter/attach_files/526851/download/tplink-adapter-0.6.2.tgz",
	//	"3143e1866673bb838297c821a253a4bf4e5fc7de2f732a10cbe7fe458fabe719", true)
	//_ = gateway.AddonsManager.InstallAddonFromUrl("yeelight-adapter",
	//	"https://gitee.com/liu_guilin/yeelight-adapter/attach_files/528153/download/yeelight-adapter-0.01.tgz",
	//	"13cce1028286b2c21afa9f9d44881e2523a2bddc64422796e6ebd86cf547e8d8", true)

	gateway.AddonsManager.LoadAddons()

}

func StartGatewayAPP(gw *HomeGateway) {
	conf := app.Config{
		HttpPort:    gw.Runtime.Ports["http"],
		HttpsPort:   gw.Runtime.Ports["https"],
		StaticDir:   "/Users/liuguilin/Documents/web-things/gateway/web-app/dist",
		TemplateDir: "/Users/liuguilin/Documents/web-things/gateway/web-app/dist",
		UploadDir:   gw.Runtime.UploadDir,
		LogDir:      gw.Runtime.LogDir,
	}
	conf.HttpPort = gw.Runtime.Ports["http"]
	webApp := app.NewApp(conf)
	webApp.Start()
}
