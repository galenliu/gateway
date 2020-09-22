package addon

//import (
//	"gateway/accessory"
//	"gateway/db"
//	"gateway/logger"
//	"io/ioutil"
//)
//
//type AddonsManager interface {
//	LoadAddons()
//}
//
//type addonsManager struct {
//	addonPath       string
//	database        db.Settings
//	InstalledAddons map[string]*Addon
//	Container       *accessory.Container
//	AddonsLoaded    bool
//}
//type addonManagerConfig struct {
//	addonPath string
//}
//
//func NewAddonsManager(addonPath string, settings db.Settings, container *accessory.Container) AddonsManager {
//	addonManager := addonsManager{
//		addonPath:       addonPath,
//		InstalledAddons: make(map[string]*Addon),
//		Container:       container,
//		database:        settings,
//		AddonsLoaded:    false,
//	}
//	addonManager.LoadAddons()
//	return &addonManager
//}
//
//func (m *addonsManager) loadAddon(addonDir string, packetId string) {
//
//	packetManifest := LoadManifest(addonDir, packetId)
//	if packetManifest == nil {
//		return
//	}
//	savedSetting := m.database.SettingWithId(packetId)
//
//	if savedSetting != nil {
//		packetManifest.Enabled = savedSetting.Enabled
//	}
//	m.database.SaveSetting(packetManifest.toSetting())
//	db := m.database.SettingWithId("buspro-adapter")
//	logger.Info.Print(db)
//
//	m.InstalledAddons[packetManifest.ID] = &Addon{Manifest: packetManifest}
//}
//
//func (m *addonsManager) LoadAddons() {
//	if m.AddonsLoaded == true {
//		return
//	}
//	fds, err := ioutil.ReadDir(m.addonPath)
//	if err != nil {
//		logger.Warning.Print("load addons err: $v", err)
//	}
//	for _, d := range fds {
//		if d.IsDir() {
//			m.loadAddon(m.addonPath, d.Name())
//		}
//	}
//	m.AddonsLoaded = true
//}
//*/
