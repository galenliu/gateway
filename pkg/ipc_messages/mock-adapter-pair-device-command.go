// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type MockAdapterPairDeviceCommandJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// Description of the device
	DeviceDescr MockAdapterPairDeviceCommandJsonDataDeviceDescr `json:"deviceDescr" yaml:"deviceDescr"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// Description of the device
type MockAdapterPairDeviceCommandJsonDataDeviceDescr map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *MockAdapterPairDeviceCommandJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in MockAdapterPairDeviceCommandJsonData: required")
	}
	if v, ok := raw["deviceDescr"]; !ok || v == nil {
		return fmt.Errorf("field deviceDescr in MockAdapterPairDeviceCommandJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in MockAdapterPairDeviceCommandJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in MockAdapterPairDeviceCommandJsonData: required")
	}
	type Plain MockAdapterPairDeviceCommandJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MockAdapterPairDeviceCommandJsonData(plain)
	return nil
}

// Tell the mock adapter to pair a device
type MockAdapterPairDeviceCommandJson struct {
	// Message-specific data
	Data MockAdapterPairDeviceCommandJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MockAdapterPairDeviceCommandJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in MockAdapterPairDeviceCommandJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in MockAdapterPairDeviceCommandJson: required")
	}
	type Plain MockAdapterPairDeviceCommandJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MockAdapterPairDeviceCommandJson(plain)
	return nil
}
