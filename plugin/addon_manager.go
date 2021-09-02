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

// GetInstallAddons 获取已安装的add-on
func (m *Manager) GetInstallAddons() []byte {
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
	err := addonInfo.setEnabled(true)

	err = m.loadAddon(addonInfo.dir, addonId)
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
	err := addonInfo.setEnabled(false)
	err = m.unloadAddon(addonId)
	if err != nil {
		return err
	}
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
		return fmt.Errorf("download addon err,pakage ID:%s err:%s", id, err.Error())
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

func (m *Manager) UninstallAddon(addonId string, disable bool) error {

	err := m.unloadAddon(addonId)
	if err != nil {
		return err
	}
	f := m.findPluginPath(addonId)
	if f != "" {
		err := util.RemoveDir(f)
		if err != nil {
			m.logger.Errorf("remove dir from: %s err :%s", f, err)
		}
	}
	if disable {
		addonInfo := m.getInstallAddon(addonId)
		err := addonInfo.setEnabled(disable)
		if err != nil {
			return err
		}
	}
	m.installAddons.Delete(addonId)
	m.installAddons.Delete(addonId)
	return nil
}

func (m *Manager) GetAddonLicense(addonId string) (string, error) {
	addonDir := m.findPluginPath(addonId)
	if addonDir == "" {
		return "", fmt.Errorf("can not find addon")
	}
	data, err := os.ReadFile(path.Join(addonDir, "LICENSE"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
