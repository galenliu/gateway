package addons

import (
	messages "gitee.com/liu_guilin/WebThings-schema"
)

type Adapter struct {
	ID          string `json:"id"`
	PackageName string `json:"package_name"`
	verbose     bool
	devices     map[string]*DeviceProxy
	manager     *Manager
	userProfile *messages.UserProfile
	preferences *messages.Preferences
}
