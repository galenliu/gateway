package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

// GetInstallAddons  获取已安装的add-on
func (m *Manager) GetInstallAddons() interface{} {
	return m.getInstallAddons()
}

func (m *Manager) EnableAddon(packageId string) error {
	addonInfo := m.getInstallAddon(packageId)
	if addonInfo == nil {
		return fmt.Errorf("package not installed")
	}
	err := addonInfo.SetEnabled(true)
	m.loadAddon(packageId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) DisableAddon(packageId string) error {
	addon := m.getInstallAddon(packageId)
	if addon == nil {
		return fmt.Errorf("package not installed")
	}
	err := addon.SetEnabled(false)
	if err != nil {
		return err
	}
	m.UnloadAddon(packageId)
	return nil
}

func (m *Manager) AddonEnabled(addonId string) bool {
	addon := m.getInstallAddon(addonId)
	return addon.Enabled
}

func (m *Manager) InstallAddonFromUrl(id, url, checksum string) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")
	m.logger.Infof("fetching add-on %s as %s", url, destPath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get err: %s", err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
		err := os.Remove(destPath)
		if err != nil {
			m.logger.Infof("remove temp file failed ,err:%s", err.Error())
		}
	}()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(destPath, data, 777)
	if !util.CheckSum(destPath, checksum) {
		return fmt.Errorf("checksum err,pakage Id:%s", id)
	}
	m.logger.Infof("download %s successful", id)
	err = m.installAddon(id, destPath)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) UninstallAddon(packetId string, disable bool) error {
	// delete addon form manager
	defer m.installAddons.Delete(packetId)
	m.UnloadAddon(packetId)
	err := util.RemoveDir(m.getAddonPath(packetId))
	if err != nil {
		m.logger.Errorf("delete plugin addon dir failed err:", err.Error())
	}
	err = util.RemoveDir(path.Join(m.config.UserProfile.DataDir, packetId))
	if err != nil {
		m.logger.Errorf("delete plugin data dir failed err:", err.Error())
	}
	return nil
}

func (m *Manager) LoadAddon(packageId string) error {
	m.loadAddon(packageId)
	return nil
}

func (m *Manager) UnloadAddon(packageId string) {
	if !m.addonsLoaded {
		m.logger.Info("The add-ons are not currently loaded, no need to unload.")
		return
	}
	plugin := m.getPlugin(packageId)
	if plugin == nil {
		m.logger.Info("The add-ons are not  register.")
	}
	plugin.unloadComponents()
	time.AfterFunc(time.Duration(constant.UnloadPluginKillDelay)*time.Millisecond, func() {
		plugin.shutdown()
		_, _ = m.extensions.LoadAndDelete(packageId)
	})
}

func (m *Manager) GetAddonLicense(addonId string) (string, error) {
	addonDir := m.getAddonPath(addonId)
	m.pluginServer.getPlugin(addonId)
	if addonDir == "" {
		return "", fmt.Errorf("addon not installed")
	}
	data, err := os.ReadFile(path.Join(addonDir, "LICENSE"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
