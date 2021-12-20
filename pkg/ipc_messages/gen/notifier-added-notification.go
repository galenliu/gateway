// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type NotifierAddedNotificationJsonData struct {
	// Name of the new notifier
	Name string `json:"name" yaml:"name"`

	// ID of the new notifier
	NotifierId string `json:"notifierId" yaml:"notifierId"`

	// Name of the add-on package
	PackageName string `json:"packageName" yaml:"packageName"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NotifierAddedNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in NotifierAddedNotificationJsonData: required")
	}
	if v, ok := raw["notifierId"]; !ok || v == nil {
		return fmt.Errorf("field notifierId in NotifierAddedNotificationJsonData: required")
	}
	if v, ok := raw["packageName"]; !ok || v == nil {
		return fmt.Errorf("field packageName in NotifierAddedNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in NotifierAddedNotificationJsonData: required")
	}
	type Plain NotifierAddedNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = NotifierAddedNotificationJsonData(plain)
	return nil
}

// Notification that a plugin has added a notifier
type NotifierAddedNotificationJson struct {
	// Message-specific data
	Data NotifierAddedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NotifierAddedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in NotifierAddedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in NotifierAddedNotificationJson: required")
	}
	type Plain NotifierAddedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = NotifierAddedNotificationJson(plain)
	return nil
}