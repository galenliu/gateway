package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"

	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type AddonManager interface {
	gateway.Component
	GetDevices() []addon.IDevice
	GetDevice(id string) addon.IDevice

	SetPIN(thingId string, pin interface{}) error

	AddNewThing(timeOut float64) error
	CancelAddNewThing()

	GetPropertyValue(deviceId string, propertyName string) (value interface{}, err error)
	SetPropertyValue(deviceId string, propertyName string, newValue interface{}) (value interface{}, err error)

	InstallAddonFromUrl(id, url, checksum string, enabled bool) error
	GetInstallAddons() ([]byte, error)
	EnableAddon(id string) error
	DisableAddon(id string) error
}

// GetInstallAddons 获取已安装的add-on
func (m *manager) GetInstallAddons() ([]byte, error) {
	m.locker.Lock()
	defer m.locker.Unlock()
	var addons []*internal.AddonInfo
	for _, v := range m.installAddons {
		addons = append(addons, v)
	}
	data, err := json.Marshal(addons)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *manager) EnableAddon(addonId string) error {
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

func (m *manager) DisableAddon(addonId string) error {
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

func (m *manager) SetPropertyValue(deviceId, propName string, newValue interface{}) (interface{}, error) {

	go func() {
		err := m.handleSetProperty(deviceId, propName, newValue)
		if err != nil {
			m.logger.Error(err.Error())
		}
	}()
	closeChan := make(chan struct{})
	propChan := make(chan interface{})
	time.AfterFunc(3*time.Second, func() {
		closeChan <- struct{}{}
	})
	changed := func(data []byte) {
		id := json.Get(data, "deviceId").ToString()
		name := json.Get(data, "name").ToString()
		value := json.Get(data, "value").GetInterface()
		if id == deviceId && name == propName {
			propChan <- value
		}
	}
	go m.bus.Subscribe(util.PropertyChanged, changed)
	defer m.bus.Unsubscribe(util.PropertyChanged, changed)
	for {
		select {
		case v := <-propChan:
			return v, nil
		case <-closeChan:
			m.logger.Error("set property(name: %s value: %s) timeout", propName, newValue)
			return nil, fmt.Errorf("timeout")
		}
	}
}

func (m *manager) GetDevice(deviceId string) addon.IDevice {
	device, ok := m.devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func (m *manager) GetDevices() (device []addon.IDevice) {
	for _, dev := range m.devices {
		device = append(device, dev)
	}
	return
}

func (m *manager) RemoveDevice(deviceId string) error {

	device := m.getDevice(deviceId)
	adapter := m.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.removeThing(device)
		return nil
	}
	return fmt.Errorf("can not find thing")
}

func (m *manager) GetPropertyValue(deviceId, propName string) (interface{}, error) {
	device, ok := m.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
	}
	prop := device.GetProperty(propName)
	if prop == nil {
		return nil, fmt.Errorf("propName(%s)invaild", propName)
	}
	return prop.GetValue(), nil
}

func (m *manager) InstallAddonFromUrl(id, url, checksum string, enabled bool) error {

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

func (m *manager) AddNewThing(pairingTimeout float64) error {
	if m.isPairing {
		return fmt.Errorf("add already in progress")
	}
	for _, adapter := range m.adapters {
		adapter.pairing(pairingTimeout)
	}
	m.isPairing = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(pairingTimeout)*time.Millisecond)
	var handlePairingTimeout = func() {
		for {
			select {
			case <-ctx.Done():
				cancelFn()
				m.CancelAddNewThing()
				//bus.Publish(util.PairingTimeout)
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func (m *manager) CancelAddNewThing() {
	if !m.isPairing {
		return
	}
	for _, adapter := range m.adapters {
		adapter.cancelPairing()
	}
	m.isPairing = false
	return
}

func (m *manager) CancelRemoveThing(deviceId string) {
	device := m.getDevice(deviceId)
	if device == nil {
		return
	}
	adapter := m.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.cancelRemoveThing(deviceId)
	}
}

func (m *manager) SetPIN(thingId string, pin interface{}) error {
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

func (m *manager) UninstallAddon(addonId string, disable bool) error {
	var key = "addons." + addonId
	err := m.unloadAddon(addonId)
	if err != nil {
		return err
	}
	f, e := m.findPlugin(addonId)
	if e == nil {
		util.RemoveDir(f)
	}
	util.RemoveDir(path.Join(path.Join(m.options.DataDir, util.DataDirName), addonId))

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
