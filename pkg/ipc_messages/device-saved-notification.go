// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Notification that the user has saved a device
type DeviceSavedNotificationJson struct {
	// Message-specific data
	Data DeviceSavedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSavedNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceSavedNotificationJsonData: required")
	}
	if v, ok := raw["device"]; !ok || v == nil {
		return fmt.Errorf("field device in DeviceSavedNotificationJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceSavedNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceSavedNotificationJsonData: required")
	}
	type Plain DeviceSavedNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSavedNotificationJsonData(plain)
	return nil
}

// Message-specific data
type DeviceSavedNotificationJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// Device corresponds to the JSON schema field "device".
	Device DeviceWithoutId `json:"device" yaml:"device"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSavedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceSavedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceSavedNotificationJson: required")
	}
	type Plain DeviceSavedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSavedNotificationJson(plain)
	return nil
}
