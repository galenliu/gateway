// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type MockAdapterRemoveDeviceResponseJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId *string `json:"deviceId,omitempty" yaml:"deviceId,omitempty"`

	// Error message in the case of failure
	Error *string `json:"error,omitempty" yaml:"error,omitempty"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Whether or not the operation was successful
	Success bool `json:"success" yaml:"success"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MockAdapterRemoveDeviceResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in MockAdapterRemoveDeviceResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in MockAdapterRemoveDeviceResponseJsonData: required")
	}
	if v, ok := raw["success"]; !ok || v == nil {
		return fmt.Errorf("field success in MockAdapterRemoveDeviceResponseJsonData: required")
	}
	type Plain MockAdapterRemoveDeviceResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MockAdapterRemoveDeviceResponseJsonData(plain)
	return nil
}

// Notice that the mock adapter has finished remove a device
type MockAdapterRemoveDeviceResponseJson struct {
	// Message-specific data
	Data MockAdapterRemoveDeviceResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MockAdapterRemoveDeviceResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in MockAdapterRemoveDeviceResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in MockAdapterRemoveDeviceResponseJson: required")
	}
	type Plain MockAdapterRemoveDeviceResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MockAdapterRemoveDeviceResponseJson(plain)
	return nil
}
