// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type OutletNotifyResponseJsonData struct {
	// ID of the request message
	MessageId int `json:"messageId" yaml:"messageId"`

	// ID of the notifier
	NotifierId string `json:"notifierId" yaml:"notifierId"`

	// ID of the outlet
	OutletId string `json:"outletId" yaml:"outletId"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Whether or not the operation was successful
	Success bool `json:"success" yaml:"success"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletNotifyResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in OutletNotifyResponseJsonData: required")
	}
	if v, ok := raw["notifierId"]; !ok || v == nil {
		return fmt.Errorf("field notifierId in OutletNotifyResponseJsonData: required")
	}
	if v, ok := raw["outletId"]; !ok || v == nil {
		return fmt.Errorf("field outletId in OutletNotifyResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in OutletNotifyResponseJsonData: required")
	}
	if v, ok := raw["success"]; !ok || v == nil {
		return fmt.Errorf("field success in OutletNotifyResponseJsonData: required")
	}
	type Plain OutletNotifyResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletNotifyResponseJsonData(plain)
	return nil
}

// Notice that an outlet notification has finished
type OutletNotifyResponseJson struct {
	// Message-specific data
	Data OutletNotifyResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OutletNotifyResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in OutletNotifyResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in OutletNotifyResponseJson: required")
	}
	type Plain OutletNotifyResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OutletNotifyResponseJson(plain)
	return nil
}
