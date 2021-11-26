// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type AdapterPairingPromptNotificationJsonData struct {
	// ID of the adapter
	AdapterId string `json:"adapterId" yaml:"adapterId"`

	// ID of specific device the prompt pertains to
	DeviceId *string `json:"deviceId,omitempty" yaml:"deviceId,omitempty"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// The message to present to the user
	Prompt string `json:"prompt" yaml:"prompt"`

	// URL of a web page containing more information
	Url *string `json:"url,omitempty" yaml:"url,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterPairingPromptNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["adapterId"]; !ok || v == nil {
		return fmt.Errorf("field adapterId in AdapterPairingPromptNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in AdapterPairingPromptNotificationJsonData: required")
	}
	if v, ok := raw["prompt"]; !ok || v == nil {
		return fmt.Errorf("field prompt in AdapterPairingPromptNotificationJsonData: required")
	}
	type Plain AdapterPairingPromptNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterPairingPromptNotificationJsonData(plain)
	return nil
}

// Notification that a prompt should be presented to the user while pairing
type AdapterPairingPromptNotificationJson struct {
	// Message-specific data
	Data AdapterPairingPromptNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AdapterPairingPromptNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in AdapterPairingPromptNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in AdapterPairingPromptNotificationJson: required")
	}
	type Plain AdapterPairingPromptNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AdapterPairingPromptNotificationJson(plain)
	return nil
}