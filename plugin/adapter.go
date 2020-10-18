package plugin

import (
	messages "gitee.com/liu_guilin/WebThings-schema"
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
