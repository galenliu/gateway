// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type OutletRemovedNotificationJsonData struct {
	// ID of the notifier
	NotifierId string `json:"notifierId" yaml:"notifierId"`

	// ID of the outlet which was removed
	OutletId string `json:"outletId" yaml:"outletId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletRemovedNotificationJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["notifierId"]; !ok || v == nil {
		return fmt.Errorf("field notifierId in OutletRemovedNotificationJsonData: required")
	}
	if v, ok := raw["outletId"]; !ok || v == nil {
		return fmt.Errorf("field outletId in OutletRemovedNotificationJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in OutletRemovedNotificationJsonData: required")
	}
	type Plain OutletRemovedNotificationJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletRemovedNotificationJsonData(plain)
	return nil
}

// Notification that an outlet has been removed from a notifier
type OutletRemovedNotificationJson struct {
	// Message-specific data
	Data OutletRemovedNotificationJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletRemovedNotificationJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in OutletRemovedNotificationJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in OutletRemovedNotificationJson: required")
	}
	type Plain OutletRemovedNotificationJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletRemovedNotificationJson(plain)
	return nil
}
