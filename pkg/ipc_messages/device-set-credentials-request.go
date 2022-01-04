// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type DeviceSetCredentialsRequestJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of the device
	DeviceId string `json:"deviceId" yaml:"deviceId"`

	// Unique ID of this message
	MessageId int `json:"messageId" yaml:"messageId"`

	// Password to set
	Password string `json:"password" yaml:"password"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Username to set
	Username string `json:"username" yaml:"username"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetCredentialsRequestJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in DeviceSetCredentialsRequestJsonData: required")
	}
	if v, ok := raw["deviceId"]; !ok || v == nil {
		return fmt.Errorf("field deviceId in DeviceSetCredentialsRequestJsonData: required")
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in DeviceSetCredentialsRequestJsonData: required")
	}
	if v, ok := raw["password"]; !ok || v == nil {
		return fmt.Errorf("field password in DeviceSetCredentialsRequestJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in DeviceSetCredentialsRequestJsonData: required")
	}
	if v, ok := raw["username"]; !ok || v == nil {
		return fmt.Errorf("field username in DeviceSetCredentialsRequestJsonData: required")
	}
	type Plain DeviceSetCredentialsRequestJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetCredentialsRequestJsonData(plain)
	return nil
}

// Set the credentials on a device
type DeviceSetCredentialsRequestJson struct {
	// Message-specific data
	Data DeviceSetCredentialsRequestJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DeviceSetCredentialsRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in DeviceSetCredentialsRequestJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in DeviceSetCredentialsRequestJson: required")
	}
	type Plain DeviceSetCredentialsRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DeviceSetCredentialsRequestJson(plain)
	return nil
}
