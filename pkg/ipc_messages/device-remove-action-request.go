// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type DeviceRemoveActionRequestJsonData struct {
	// ID of the existing actions
	ActionId string `json:"actionId" yaml:"actionId"`

	// Name of the actions
	ActionName string `json:"actionName" yaml:"actionName"`

	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// Unique ID of this message
	MessageId int `json:"messageId" yaml:"messageId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRemoveActionRequestJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["actionId"]; !ok || v == nil {
		return fmt.Errorf("field actionId in DeviceRemoveActionRequestJsonData: required")
	}
	if v, ok := raw["actionName"]; !ok || v == nil {
		return fmt.Errorf("field actionName in DeviceRemoveActionRequestJsonData: required")
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceRemoveActionRequestJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceRemoveActionRequestJsonData: required")
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in DeviceRemoveActionRequestJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceRemoveActionRequestJsonData: required")
	}
	type Plain DeviceRemoveActionRequestJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRemoveActionRequestJsonData(plain)
	return nil
}

// Remove/cancel an existing actions from a device
type DeviceRemoveActionRequestJson struct {
	// Message-specific data
	Data DeviceRemoveActionRequestJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRemoveActionRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceRemoveActionRequestJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceRemoveActionRequestJson: required")
	}
	type Plain DeviceRemoveActionRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRemoveActionRequestJson(plain)
	return nil
}
