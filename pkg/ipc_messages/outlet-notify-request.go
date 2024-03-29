// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletNotifyRequestJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["level"]; !ok || v == nil {
		return fmt.Errorf("field level in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["message"]; !ok || v == nil {
		return fmt.Errorf("field message in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["notifierId"]; !ok || v == nil {
		return fmt.Errorf("field notifierId in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["outletId"]; !ok || v == nil {
		return fmt.Errorf("field outletId in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in OutletNotifyRequestJsonData: required")
	}
	if v, ok := raw["title"]; !ok || v == nil {
		return fmt.Errorf("field title in OutletNotifyRequestJsonData: required")
	}
	type Plain OutletNotifyRequestJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletNotifyRequestJsonData(plain)
	return nil
}

// Notify a user via an outlet
type OutletNotifyRequestJson struct {
	// Message-specific data
	Data OutletNotifyRequestJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// Message-specific data
type OutletNotifyRequestJsonData struct {
	// Priority level of the notification, 0 being the lowest priority
	Level NotificationLevel `json:"level" yaml:"level"`

	// Message of the notification
	Message string `json:"message" yaml:"message"`

	// Unique ID of this message
	MessageId int `json:"messageId" yaml:"messageId"`

	// ID of the notifier
	NotifierId string `json:"notifierId" yaml:"notifierId"`

	// ID of the outlet
	OutletId string `json:"outletId" yaml:"outletId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Title of the notification
	Title string `json:"title" yaml:"title"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletNotifyRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in OutletNotifyRequestJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in OutletNotifyRequestJson: required")
	}
	type Plain OutletNotifyRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletNotifyRequestJson(plain)
	return nil
}
