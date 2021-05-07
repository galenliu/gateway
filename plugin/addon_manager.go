package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"path"
	"sync"
	"time"
)

type Map = map[string]interface{}

var instance *AddonManager

func NewAddonsManager() *AddonManager {
	am := &AddonManager{}
	instance = am

	am.AddonsDir = configs.GetAddonsDir()
	am.DataDir = configs.GetAddonsDir()
	am.verbose = configs.IsVerbose()
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.devices = make(map[string]addon.IDevice, 50)
	am.installAddons = make(map[string]*internal.AddonInfo, 50)
	am.adapters = make(map[string]*Adapter, 20)
	am.extensions = make(map[string]Extension)

	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

//获取已安装的add-on
func GetInstallAddons() ([]byte, error) {
	instance.locker.Lock()
	defer instance.locker.Unlock()
	var addons []*internal.AddonInfo
	for _, v := range instance.installAddons {
		addons = append(addons, v)
	}
	data, err := json.Marshal(addons)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EnableAddon(addonId string) error {
	instance.locker.Lock()
	defer instance.locker.Unlock()
	addonInfo := instance.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not exit")
	}

	err := addonInfo.UpdateAddonInfoToDB(true)
	err = instance.loadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func DisableAddon(addonId string) error {
	instance.locker.Lock()
	defer instance.locker.Unlock()
	addonInfo := instance.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.UpdateAddonInfoToDB(false)
	if err != nil {
		return err
	}
	err = instance.unloadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func GetDevices() (device []addon.IDevice) {
	for _, dev := range instance.devices {
		device = append(device, dev)
	}
	return
}

func GetDevice(deviceId string) addon.IDevice {
	device, ok := instance.devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func SetProperty(deviceId, propName string, newValue interface{}) (interface{}, error) {

	go func() {
		err := instance.handleSetProperty(deviceId, propName, newValue)
		if err != nil {
			log.Error(err.Error())
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
	go Subscribe(util.PropertyChanged, changed)
	defer Unsubscribe(util.PropertyChanged, changed)
	for {
		select {
		case v := <-propChan:
			return v, nil
		case <-closeChan:
			log.Error("set property(name: %s value: %s) timeout", propName, newValue)
			return nil, fmt.Errorf("timeout")
		}
	}
}

func RemoveDevice(deviceId string) error {

	device := instance.getDevice(deviceId)
	adapter := instance.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.removeThing(device)
		return nil
	}
	return fmt.Errorf("can not find thing")
}

func GetPropertyValue(deviceId, propName string) (interface{}, error) {
	device, ok := instance.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
	}
	prop := device.GetProperty(propName)
	if prop == nil {
		return nil, fmt.Errorf("propName(%s)invaild", propName)
	}
	return prop.GetValue(), nil
}

func InstallAddonFromUrl(id, url, checksum string, enabled bool) error {
	return instance.installAddonFromUrl(id, url, checksum, enabled)
}

func AddNewThing(pairingTimeout float64) error {
	if instance.isPairing {
		return fmt.Errorf("add already in progress")
	}
	for _, adapter := range instance.adapters {
		adapter.pairing(pairingTimeout)
	}
	instance.isPairing = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(pairingTimeout)*time.Millisecond)
	var handlePairingTimeout = func() {
		for {
			select {
			case <-ctx.Done():
				cancelFn()
				CancelAddNewThing()
				//bus.Publish(util.PairingTimeout)
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func CancelAddNewThing() {
	if !instance.isPairing {
		return
	}
	for _, adapter := range instance.adapters {
		adapter.cancelPairing()
	}
	instance.isPairing = false
	return
}

func CancelRemoveThing(deviceId string) {
	device := instance.getDevice(deviceId)
	if device == nil {
		return
	}
	adapter := instance.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.cancelRemoveThing(deviceId)
	}
}

func SetThingPin(thingId string, pin interface{}) error {
	device := instance.getDevice(thingId)
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

func RequestAction(thingId, actionId, actionName string, actionParams interface{}) error {
	//TODO
	return nil
}

func UnloadAddon(addonId string) error {
	return instance.unloadAddon(addonId)
}
func LoadAddon(addonId string) error {
	return instance.loadAddon(addonId)
}

func AddonEnabled(addonId string) bool {
	a, ok := instance.installAddons[addonId]
	if !ok {
		return false
	}
	return a.Enabled
}

func UninstallAddon(addonId string, disable bool) error {
	var key = "addons." + addonId
	err := instance.unloadAddon(addonId)
	if err != nil {
		return err
	}
	util.RemoveDir(path.Join(instance.DataDir, addonId))
	util.RemoveDir(path.Join(instance.AddonsDir, addonId))

	if disable {
		setting, err := database.GetSetting(key)
		if err != nil {
			log.Error(err.Error())
		}
		var addonInfo internal.AddonInfo
		err = json.UnmarshalFromString(setting, &addonInfo)
		_ = addonInfo.UpdateAddonInfoToDB(false)
	}
	delete(instance.installAddons, addonId)
	return nil
}

func Subscribe(typ string, f interface{}) {
	_ = bus.Subscribe(typ, f)
}

func Unsubscribe(typ string, f interface{}) {
	_ = bus.Unsubscribe(typ, f)
}

func Publish(typ string, args ...interface{}) {
	bus.Publish(typ, args...)
}
