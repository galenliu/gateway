// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Notification that a property on a device has changed
type DevicePropertyChangedNotificationJson struct {
	// Message-specific data
	Data DevicePropertyChangedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DevicePropertyChangedNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DevicePropertyChangedNotificationJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DevicePropertyChangedNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DevicePropertyChangedNotificationJsonData: required")
	}
	if v, ok := raw["property"]; !ok || v == nil {
		return fmt.Errorf("field property in DevicePropertyChangedNotificationJsonData: required")
	}
	type Plain DevicePropertyChangedNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DevicePropertyChangedNotificationJsonData(plain)
	return nil
}

// Message-specific data
type DevicePropertyChangedNotificationJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Property corresponds to the JSON schema field "property".
	Property Property `json:"property" yaml:"property"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DevicePropertyChangedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DevicePropertyChangedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DevicePropertyChangedNotificationJson: required")
	}
	type Plain DevicePropertyChangedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DevicePropertyChangedNotificationJson(plain)
	return nil
}
