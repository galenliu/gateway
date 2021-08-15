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
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.enable()

	err = m.loadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) DisableAddon(addonId string) error {
	addonInfo := m.getInstallAddon(addonId)
	if addonInfo == nil {
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.disable()
	err = m.unloadAddon(addonId)
	if err != nil {
		return err
	}
	err = m.loadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) InstallAddonFromUrl(id, url, checksum string, enabled bool) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")
	m.logger.Info("fetching add-on %s as %s", url, destPath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Download addon err,pakage ID:%s err:%s", id, err.Error()))
	}
	defer func() {
		_ = resp.Body.Close()
		err := os.Remove(destPath)
		if err != nil {
			m.logger.Info("remove temp file failed ,err:%s", err.Error())
		}
	}()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(destPath, data, 777)
	if !util.CheckSum(destPath, checksum) {
		return fmt.Errorf(fmt.Sprintf("checksum err,pakage ID:%s", id))
	}
	err = m.installAddon(id, destPath, enabled)
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
			m.logger.Error("remove dir from: %s err :%s", f, err)
		}
	}
	if disable {
		addonInfo := m.getInstallAddon(addonId)
		err := addonInfo.disable()
		if err != nil {
			return err
		}
	}
	m.installAddons.Delete(addonId)
	m.installAddons.Delete(addonId)
	return nil
}

func (m *Manager) GetAddonLicense(addonId string) string {
	return ""
}

//func Subscribe(typ string, f interface{}) {
//	_ = event_bus.Subscribe(typ, f)
//}
//
//func Unsubscribe(typ string, f interface{}) {
//	_ = event_bus.Unsubscribe(typ, f)
//}
//
//func Publish(typ string, args ...interface{}) {
//	event_bus.Publish(typ, args...)
//}
