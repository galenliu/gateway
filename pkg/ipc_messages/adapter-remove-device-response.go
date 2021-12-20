// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type AdapterRemoveDeviceResponseJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device which has been removed
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterRemoveDeviceResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in AdapterRemoveDeviceResponseJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in AdapterRemoveDeviceResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in AdapterRemoveDeviceResponseJsonData: required")
	}
	type Plain AdapterRemoveDeviceResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterRemoveDeviceResponseJsonData(plain)
	return nil
}

// Notice that an adapter has removed a device
type AdapterRemoveDeviceResponseJson struct {
	// Message-specific data
	Data AdapterRemoveDeviceResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterRemoveDeviceResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in AdapterRemoveDeviceResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in AdapterRemoveDeviceResponseJson: required")
	}
	type Plain AdapterRemoveDeviceResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterRemoveDeviceResponseJson(plain)
	return nil
}
