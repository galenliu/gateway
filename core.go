package gateway

import (
	"gateway/db"
	"gateway/plugin"
	messages "gitee.com/liu_guilin/WebThings-schema"
	"path"
)

type HomeGateway struct {
	UserProfile   messages.UserProfile
	Preferences   messages.Preferences
	AddonsManager *plugin.AddonsManager
}

func NewHomeGateway(baseDir string) *HomeGateway {
	var gateway = &HomeGateway{}
	gateway.setUserProfile(baseDir)
	return gateway
}

func (gateway *HomeGateway) Run() error {
	gateway.AddonsManager.Run()
	return nil
}

func (gateway *HomeGateway) updatePreferences() error {
	lang := "preferences.language"
	temp := "preferences.units.temperature"
	l, err := db.DB.SettingGet(lang)
	if l == "" && err != nil {
		l = "zh-CN"
		err := db.DB.SettingSet(lang, l)
		if err != nil {
			return err
		}
	}
	var t string
	t, err = db.DB.SettingGet(temp)
	if t == "" && err != nil {
		t = "degree celsius"
		err = db.DB.SettingSet(temp, t)
		if err != nil {
			return err
		}
	}
	unit := messages.Units{Temperature: t}
	preferences := messages.Preferences{
		Language: l,
		Units:    unit,
	}
	gateway.Preferences = preferences
	return nil
}

func (gateway *HomeGateway) setUserProfile(baseDir string) {
	var profile = messages.UserProfile{
		BaseDir:        baseDir,
		DataDir:        path.Join(baseDir, DataDir),
		AddonsDir:      path.Join(baseDir, AddonsDir),
		ConfigDir:      path.Join(baseDir, ConfigDir),
		MediaDir:       path.Join(baseDir, MediaDir),
		UploadDir:      path.Join(baseDir, UploadDir),
		LogDir:         path.Join(baseDir, LogDir),
		GatewayVersion: Version,
	}
	EnsureConfigPath(profile.BaseDir, profile.DataDir, profile.AddonsDir, profile.ConfigDir, profile.MediaDir, profile.UploadDir, profile.LogDir)
	gateway.UserProfile = profile
}

func (gateway *HomeGateway) addonManagerLoadAndRun() error {
	addonManager := plugin.NewAddonsManager(gateway, Log)
	gateway.AddonsManager = addonManager
	addonManager.LoadAddons()
	return nil
}

func (gateway *HomeGateway) Close() {

}

func (gateway *HomeGateway) GetUserProfile() *messages.UserProfile {
	return &gateway.UserProfile
}
func (gateway *HomeGateway) GetPreferences() *messages.Preferences {
	return &gateway.Preferences
}

func (gateway *HomeGateway) EnsureConfigPath(dir string, dirs ...string) {
	EnsureConfigPath(dir)
}
