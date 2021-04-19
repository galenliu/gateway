package plugin

import (
	"addon"
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
	"path"
	"time"
)

type Map = map[string]interface{}

func GetInstallAddons() ([]byte, error) {
	instance.locker.Lock()
	defer instance.locker.Unlock()
	var addons []*AddonInfo
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

	go instance.handleSetProperty(deviceId, propName, newValue)
	closeChan := make(chan struct{})
	propChan := make(chan Map)
	time.AfterFunc(3*time.Second, func() {
		closeChan <- struct{}{}
	})
	changed := func(data Map) {
		if to.String(data["deviceId"]) == deviceId && to.String(data["name"]) == propName {
			propChan <- data
		}
	}
	go Subscribe(util.PropertyChanged, changed)
	defer Unsubscribe(util.PropertyChanged, changed)
	for {
		select {
		case data := <-propChan:
			return data["value"], nil
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
		config, err := database.GetSetting(key)
		if err != nil {
			log.Error(err.Error())
		}
		var addonInfo AddonInfo
		err = json.UnmarshalFromString(config, &addonInfo)
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
