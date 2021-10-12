package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// GetInstallAddonsBytes  获取已安装的add-on
func (m *Manager) GetInstallAddonsBytes() []byte {
	addons := m.getInstallAddons()
	data, err := json.Marshal(addons)
	if err != nil {
		return nil
	}
	return data
}

func (m *Manager) EnableAddon(addonId string) error {
	addonInfo := m.getInstallAddon(addonId)
	if addonInfo == nil {
		return fmt.Errorf("package not installed")
	}
	err := addonInfo.SetEnabled(true)

	err = m.loadAddon(addonInfo.Dir, addonId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) DisableAddon(addonId string) error {
	addonInfo := m.getInstallAddon(addonId)
	if addonInfo == nil {
		return fmt.Errorf("package not installed")
	}
	err := addonInfo.SetEnabled(false)
	if err != nil {
		return err
	}
	plugin:=m.pluginServer.findPlugin(addonId)
	if plugin==nil{
		return nil
	}
	plugin.disable()
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
		return fmt.Errorf("checksum err,pakage ID:%s", id)
	}
	m.logger.Infof("download %s successful", id)
	err = m.installAddon(id, destPath)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) UninstallAddon(pluginId string, disable bool) error {

	defer m.installAddons.Delete(pluginId)
	err := m.pluginServer.unloadPlugin(pluginId)
	if err != nil {
		return err
	}
	if disable {
		addonInfo := m.getInstallAddon(pluginId)
		err := addonInfo.SetEnabled(disable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) GetAddonLicense(addonId string) (string, error) {
	addonDir := m.findAddon(addonId)
	m.pluginServer.findPlugin(addonId)
	if addonDir == "" {
		return "", fmt.Errorf("addon not installed")
	}
	data, err := os.ReadFile(path.Join(addonDir, "LICENSE"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Manager) UnloadAddon(id string) error {
	plugin := m.pluginServer.findPlugin(id)
	if plugin == nil {
		return fmt.Errorf("plugin not exist")
	}
	return m.pluginServer.unloadPlugin(id)
}

func (m *Manager) LoadAddon(id string) error {
	addon := m.findAddon(id)
	if addon == "" {
		return fmt.Errorf("addon not installed")
	}
	return m.loadAddon(id, addon)
}