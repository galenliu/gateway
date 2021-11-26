// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Set a property value on a device
type DeviceSetPropertyCommandJson struct {
	// Message-specific data
	Data DeviceSetPropertyCommandJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetPropertyCommandJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceSetPropertyCommandJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceSetPropertyCommandJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceSetPropertyCommandJsonData: required")
	}
	if v, ok := raw["propertyName"]; !ok || v == nil {
		return fmt.Errorf("field propertyName in DeviceSetPropertyCommandJsonData: required")
	}
	if v, ok := raw["propertyValue"]; !ok || v == nil {
		return fmt.Errorf("field propertyValue in DeviceSetPropertyCommandJsonData: required")
	}
	type Plain DeviceSetPropertyCommandJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetPropertyCommandJsonData(plain)
	return nil
}

// Message-specific data
type DeviceSetPropertyCommandJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Name of the property to set
	PropertyName string `json:"propertyName" yaml:"propertyName"`

	// New value of the property
	PropertyValue DeviceSetPropertyCommandJsonDataPropertyValue `json:"propertyValue" yaml:"propertyValue"`
}

// New value of the property
type DeviceSetPropertyCommandJsonDataPropertyValue interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetPropertyCommandJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceSetPropertyCommandJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceSetPropertyCommandJson: required")
	}
	type Plain DeviceSetPropertyCommandJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetPropertyCommandJson(plain)
	return nil
}
