// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type ApiHandlerAddedNotificationJsonData struct {
	// Instance of the add-on package
	PackageName string `json:"packageName" yaml:"packageName"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerAddedNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["packageName"]; !ok || v == nil {
		return fmt.Errorf("field packageName in ApiHandlerAddedNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in ApiHandlerAddedNotificationJsonData: required")
	}
	type Plain ApiHandlerAddedNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerAddedNotificationJsonData(plain)
	return nil
}

// Notification that a plugin has added an API handler
type ApiHandlerAddedNotificationJson struct {
	// Message-specific data
	Data ApiHandlerAddedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerAddedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in ApiHandlerAddedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in ApiHandlerAddedNotificationJson: required")
	}
	type Plain ApiHandlerAddedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerAddedNotificationJson(plain)
	return nil
}
