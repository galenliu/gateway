// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Notification that an event has occurred on a device
type DeviceEventNotificationJson struct {
	// Message-specific data
	Data DeviceEventNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// Message-specific data
type DeviceEventNotificationJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// Description of the event
	Event DeviceEventNotificationJsonDataEvent `json:"event" yaml:"event"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// Description of the event
type DeviceEventNotificationJsonDataEvent struct {
	// Data corresponds to the JSON schema field "data".
	Data DeviceEventNotificationJsonDataEventData `json:"data,omitempty" yaml:"data,omitempty"`

	// Name of the event
	Name string `json:"name" yaml:"name"`

	// Timestamp of the event
	Timestamp string `json:"timestamp" yaml:"timestamp"`
}

type DeviceEventNotificationJsonDataEventData interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceEventNotificationJsonDataEvent) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in DeviceEventNotificationJsonDataEvent: required")
	}
	if v, ok := raw["timestamp"]; !ok || v == nil {
		return fmt.Errorf("field timestamp in DeviceEventNotificationJsonDataEvent: required")
	}
	type Plain DeviceEventNotificationJsonDataEvent
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceEventNotificationJsonDataEvent(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceEventNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceEventNotificationJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceEventNotificationJsonData: required")
	}
	if v, ok := raw["event"]; !ok || v == nil {
		return fmt.Errorf("field event in DeviceEventNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceEventNotificationJsonData: required")
	}
	type Plain DeviceEventNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceEventNotificationJsonData(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceEventNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceEventNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceEventNotificationJson: required")
	}
	type Plain DeviceEventNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceEventNotificationJson(plain)
	return nil
}