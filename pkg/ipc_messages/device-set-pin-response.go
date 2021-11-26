// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Notice that setting the PIN on a device has finished
type DeviceSetPinResponseJson struct {
	// Message-specific data
	Data DeviceSetPinResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetPinResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceSetPinResponseJsonData: required")
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in DeviceSetPinResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceSetPinResponseJsonData: required")
	}
	if v, ok := raw["success"]; !ok || v == nil {
		return fmt.Errorf("field success in DeviceSetPinResponseJsonData: required")
	}
	type Plain DeviceSetPinResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetPinResponseJsonData(plain)
	return nil
}

// Message-specific data
type DeviceSetPinResponseJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// Device corresponds to the JSON schema field "device".
	Device *Device `json:"device,omitempty" yaml:"device,omitempty"`

	// ID of the device
	DeviceId *string `json:"deviceId,omitempty" yaml:"deviceId,omitempty"`

	// ID of the request message
	MessageId int `json:"messageId" yaml:"messageId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Whether or not the operation was successful
	Success bool `json:"success" yaml:"success"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetPinResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceSetPinResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceSetPinResponseJson: required")
	}
	type Plain DeviceSetPinResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetPinResponseJson(plain)
	return nil
}