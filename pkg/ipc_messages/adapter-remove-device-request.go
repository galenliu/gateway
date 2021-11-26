// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type AdapterRemoveDeviceRequestJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device to remove
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterRemoveDeviceRequestJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in AdapterRemoveDeviceRequestJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in AdapterRemoveDeviceRequestJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in AdapterRemoveDeviceRequestJsonData: required")
	}
	type Plain AdapterRemoveDeviceRequestJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterRemoveDeviceRequestJsonData(plain)
	return nil
}

// Ask an adapter to remove a device
type AdapterRemoveDeviceRequestJson struct {
	// Message-specific data
	Data AdapterRemoveDeviceRequestJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterRemoveDeviceRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in AdapterRemoveDeviceRequestJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in AdapterRemoveDeviceRequestJson: required")
	}
	type Plain AdapterRemoveDeviceRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterRemoveDeviceRequestJson(plain)
	return nil
}
