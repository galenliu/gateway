package plugin

import (
	"addon"
	"context"
	"fmt"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models/thing"
	json "github.com/json-iterator/go"
	"path"
	"time"
)

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

func GetDevices() []*addon.Device {
	var devs []*addon.Device
	for _, d := range instance.devices {
		devs = append(devs, d)
		devs = append(devs, d)
	}
	return devs
}

func GetThings() []*thing.Thing {
	var ts []*thing.Thing
	for _, d := range instance.devices {
		var t = asWebThing(d)
		ts = append(ts, t)
	}
	return ts
}

func GetDevice(deviceId string) *addon.Device {
	device, ok := instance.devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func SetProperty(deviceId, propName string, newValue interface{}, ctx context.Context) (*addon.Property, error) {
	return instance.handleSetProperty(deviceId, propName, newValue, ctx)

}

func RemoveDevice(deviceId string) error {
	adapter := instance.getAdapterByDeviceId(deviceId)
	device := instance.getDevice(deviceId)
	if adapter != nil {
		adapter.removeThing(device)
		return nil
	}
	return fmt.Errorf("can not find thing")
}

func GetPropertyValue(deviceId, propName string) (interface{}, error) {
	device, ok := instance.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("devices(%s) not found", deviceId)
	}
	prop, err := device.FindProperty(propName)
	if err != nil {
		return nil, err
	}
	return prop, nil
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
	dev := instance.getDevice(deviceId)
	if dev == nil {
		return
	}
	adapter := instance.getAdapter(dev.ID)
	if adapter != nil {
		adapter.cancelRemoveThing(dev.ID)
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
