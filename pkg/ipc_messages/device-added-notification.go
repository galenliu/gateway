// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"

import "encoding/json"

// Notification that a new device has been added to an adapter
type DeviceAddedNotificationJson struct {
	// Message-specific data
	Data DeviceAddedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// Message-specific data
type DeviceAddedNotificationJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// Device corresponds to the JSON schema field "device".
	Device Device `json:"device" yaml:"device"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceAddedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceAddedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceAddedNotificationJson: required")
	}
	type Plain DeviceAddedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceAddedNotificationJson(plain)
	return nil
}