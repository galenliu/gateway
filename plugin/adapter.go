package plugin

import (
	messages "github.com/galenliu/smartassistant-ipc"
)

type Adapter struct {
	ID          string `json:"id"`
	PackageName string `json:"package_name"`
	verbose     bool
	devices     map[string]*DeviceProxy
	manager     *AddonsManager
	userProfile *messages.UserProfile
	preferences *messages.Preferences
}
