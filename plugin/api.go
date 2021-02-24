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

func FindDevice(deviceId string) (*addon.Device, error) {
	device, ok := manager.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("devices(%s) not found", deviceId)
	}
	return device, nil
}

func SetProperValue(deviceId, propName string, newValue interface{}) error {
	manager.handleSetPropertyValue(deviceId, propName, newValue)
	return nil
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

func InstallAddonFromUrl(id, url, checksum string, enabled bool) {
	_ = manager.installAddonFromUrl(id, url, checksum, enabled)
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

func SetThingPin(thingId string, pin interface{}) error {
	device, err := manager.findDevice(thingId)
	if err != nil {
		return err
	}
	err = device.SetPin(pin)
	if err != nil {
		return err
	}
	return nil
}

//func SubscriptionDeviceAdded(key interface{}, f func(*addon.Device)) func() {
//	manager.onDeviceAddedFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onDeviceAddedFuncs, key)
//	}
//	return removeFunc
//}
//
//func SubscriptionDeviceConnected(key interface{}, f func(*addon.Device, bool)) func() {
//	manager.onDeviceConnectedFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onPropertyChangedFuncs, key)
//	}
//	return removeFunc
//}
//
//func SubscriptionActionUpdate(key interface{}, f func(*addon.Action)) func() {
//	manager.onActionUpdateFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onPropertyChangedFuncs, key)
//	}
//	return removeFunc
//}
//
//func SubscriptionEvent(key interface{}, f func(*addon.Event)) func() {
//	manager.onEventFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onPropertyChangedFuncs, key)
//	}
//	return removeFunc
//}
//
//func SubscriptionDeviceRemoved(key interface{}, f func(*addon.Device)) func() {
//	manager.onDeviceRemovedFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onPropertyChangedFuncs, key)
//	}
//	return removeFunc
//}
//
//func SubscriptionPropertyChanged(key interface{}, f func(*addon.Property)) func() {
//	manager.onPropertyChangedFuncs[key] = f
//	var removeFunc = func() {
//		delete(manager.onPropertyChangedFuncs, key)
//	}
//	return removeFunc
//}
