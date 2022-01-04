// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Request a new actions on a device
type DeviceRequestActionRequestJson struct {
	// Message-specific data
	Data DeviceRequestActionRequestJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRequestActionRequestJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["actionId"]; !ok || v == nil {
		return fmt.Errorf("field actionId in DeviceRequestActionRequestJsonData: required")
	}
	if v, ok := raw["actionName"]; !ok || v == nil {
		return fmt.Errorf("field actionName in DeviceRequestActionRequestJsonData: required")
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceRequestActionRequestJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceRequestActionRequestJsonData: required")
	}
	if v, ok := raw["input"]; !ok || v == nil {
		return fmt.Errorf("field input in DeviceRequestActionRequestJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceRequestActionRequestJsonData: required")
	}
	type Plain DeviceRequestActionRequestJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRequestActionRequestJsonData(plain)
	return nil
}

// Message-specific data
type DeviceRequestActionRequestJsonData struct {
	// Unique ID of this existing actions
	ActionId string `json:"actionId" yaml:"actionId"`

	// Name of the actions
	ActionName string `json:"actionName" yaml:"actionName"`

	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// Input to the actions
	Input DeviceRequestActionRequestJsonDataInput `json:"input" yaml:"input"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// Input to the actions
type DeviceRequestActionRequestJsonDataInput any

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRequestActionRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceRequestActionRequestJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceRequestActionRequestJson: required")
	}
	type Plain DeviceRequestActionRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRequestActionRequestJson(plain)
	return nil
}
