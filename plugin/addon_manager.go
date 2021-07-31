package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"

	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)



// GetInstallAddons 获取已安装的add-on
func (m *Manager) GetInstallAddons() []byte {
	m.locker.Lock()
	defer m.locker.Unlock()
	var addons []*internal.AddonInfo
	for _, v := range m.installAddons {
		addons = append(addons, v)
	}
	data, err := json.Marshal(addons)
	if err != nil {
		return nil
	}
	return data
}

func (m *Manager) EnableAddon(addonId string) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	addonInfo := m.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not exit")
	}
	err := addonInfo.UpdateAddonInfoToDB(true)
	err = m.loadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) DisableAddon(addonId string) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	addonInfo := m.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.UpdateAddonInfoToDB(false)
	if err != nil {
		return err
	}
	err = m.unloadAddon(addonId)
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
		return fmt.Errorf(fmt.Sprintf("Download addon err,pakage id:%s err:%s", id, err.Error()))
	}
	defer func() {
		_ = resp.Body.Close()
		err := os.Remove(destPath)
		if err != nil {
			logging.Info("remove temp file failed ,err:%s", err.Error())
		}
	}()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(destPath, data, 777)
	if !util.CheckSum(destPath, checksum) {
		return fmt.Errorf(fmt.Sprintf("checksum err,pakage id:%s", id))
	}
	err = m.installAddon(id, destPath, enabled)
	if err != nil {
		return err
	}
	return nil

}

func (m *Manager) UninstallAddon(addonId string, disable bool) error {
	var key = "addons." + addonId
	err := m.unloadAddon(addonId)
	if err != nil {
		return err
	}
	f, e := m.findPlugin(addonId)
	if e == nil {
		util.RemoveDir(f)
	}
	util.RemoveDir(path.Join(path.Join(m.options.DataDir, constant.DataDirName), addonId))

	if disable {
		setting, err := database.GetSetting(key)
		if err != nil {
			logging.Error(err.Error())
		}
		var addonInfo internal.AddonInfo
		err = json.UnmarshalFromString(setting, &addonInfo)
		_ = addonInfo.UpdateAddonInfoToDB(false)
	}
	delete(m.installAddons, addonId)
	return nil
}

func (m *Manager) CancelRemoveThing(deviceId string) {
	device := m.getDevice(deviceId)
	if device == nil {
		return
	}
	adapter := m.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.cancelRemoveThing(deviceId)
	}
}

func (m *Manager) SetPIN(thingId string, pin interface{}) error {
	device := m.getDevice(thingId)
	if device == nil {
		return fmt.Errorf("con not finid device:" + thingId)
	}
	err := device.SetPin(pin)
	if err != nil {
		return err
	}
	return nil
}

func RemoveAction(thingId, actionId, actionName string) error {
	//TODO
	return nil
}

func RequestAction(thingId, actionId, actionName string, actionParams map[string]interface{}) error {
	//TODO
	return nil
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
