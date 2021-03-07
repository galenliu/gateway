package plugin

import (
	"addon"
	"context"
	"fmt"
	"gateway/server/models/thing"
	json "github.com/json-iterator/go"
	"time"
)

func GetInstallAddons() ([]byte, error) {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	var addons []*AddonInfo
	for _, v := range manager.installAddons {
		addons = append(addons, v)
	}
	data, err := json.Marshal(addons)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EnableAddon(addonId string) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	addonInfo := manager.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not exit")
	}

	err := addonInfo.UpdateAddonInfoToDB(true)
	err = manager.loadAddon(addonId, true)
	if err != nil {
		return err
	}
	return nil
}

func DisableAddon(addonId string) error {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	addonInfo := manager.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.UpdateAddonInfoToDB(false)
	if err != nil {
		return err
	}
	manager.unloadAddon(addonId)
	return nil
}

func GetDevices() []*addon.Device {
	var devs []*addon.Device
	for _, d := range manager.devices {
		devs = append(devs, d)
		devs = append(devs, d)
	}
	return devs
}

func GetThings() []*thing.Thing {
	var ts []*thing.Thing
	for _, d := range manager.devices {
		var t = asThing(d)
		ts = append(ts, t)
	}
	return ts
}

func GetDevice(deviceId string) *addon.Device {
	device, ok := manager.devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func SetPropertyValue(deviceId, propName string, newValue interface{}) error {
	return manager.handleSetPropertyValue(deviceId, propName, newValue)

}

func RemoveDevice(deviceId string) error {
	adapter := manager.getAdapterByDeviceId(deviceId)
	device := manager.getDevice(deviceId)
	if adapter != nil {
		adapter.removeThing(device)
		return nil
	}
	return fmt.Errorf("can not find thing")
}

func GetPropertyValue(deviceId, propName string) (interface{}, error) {
	device, ok := manager.devices[deviceId]
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
	return manager.installAddonFromUrl(id, url, checksum, enabled)
}

func AddNewThing(pairingTimeout float64) error {
	if manager.isPairing {
		return fmt.Errorf("add already in progress")
	}
	for _, adapter := range manager.adapters {
		adapter.pairing(pairingTimeout)
	}
	manager.isPairing = true
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
	if !manager.isPairing {
		return
	}
	for _, adapter := range manager.adapters {
		adapter.cancelPairing()
	}
	manager.isPairing = false
	return
}

func CancelRemoveThing(deviceId string) {
	dev := manager.getDevice(deviceId)
	if dev == nil {
		return
	}
	adapter := manager.getAdapter(dev.ID)
	if adapter != nil {
		adapter.cancelRemoveThing(dev.ID)
	}
}

func SetThingPin(thingId string, pin interface{}) error {
	device := manager.getDevice(thingId)
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
