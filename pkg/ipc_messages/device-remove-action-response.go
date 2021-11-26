// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type DeviceRemoveActionResponseJsonData struct {
	// ID of the action
	ActionId string `json:"actionId" yaml:"actionId"`

	// Name of the action
	ActionName string `json:"actionName" yaml:"actionName"`

	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the request message
	MessageId int `json:"messageId" yaml:"messageId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Whether or not the operation was successful
	Success bool `json:"success" yaml:"success"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRemoveActionResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["actionId"]; !ok || v == nil {
		return fmt.Errorf("field actionId in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["actionName"]; !ok || v == nil {
		return fmt.Errorf("field actionName in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceRemoveActionResponseJsonData: required")
	}
	if v, ok := raw["success"]; !ok || v == nil {
		return fmt.Errorf("field success in DeviceRemoveActionResponseJsonData: required")
	}
	type Plain DeviceRemoveActionResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRemoveActionResponseJsonData(plain)
	return nil
}

// Notice that an action has been removed/cancelled from a device
type DeviceRemoveActionResponseJson struct {
	// Message-specific data
	Data DeviceRemoveActionResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceRemoveActionResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceRemoveActionResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceRemoveActionResponseJson: required")
	}
	type Plain DeviceRemoveActionResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceRemoveActionResponseJson(plain)
	return nil
}
