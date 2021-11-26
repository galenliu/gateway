// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type AdapterCancelRemoveDeviceCommandJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device which is being removed
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterCancelRemoveDeviceCommandJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in AdapterCancelRemoveDeviceCommandJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in AdapterCancelRemoveDeviceCommandJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in AdapterCancelRemoveDeviceCommandJsonData: required")
	}
	type Plain AdapterCancelRemoveDeviceCommandJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterCancelRemoveDeviceCommandJsonData(plain)
	return nil
}

// Tell an adapter to cancel the removal of a device
type AdapterCancelRemoveDeviceCommandJson struct {
	// Message-specific data
	Data AdapterCancelRemoveDeviceCommandJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterCancelRemoveDeviceCommandJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in AdapterCancelRemoveDeviceCommandJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in AdapterCancelRemoveDeviceCommandJson: required")
	}
	type Plain AdapterCancelRemoveDeviceCommandJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterCancelRemoveDeviceCommandJson(plain)
	return nil
}
