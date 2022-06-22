// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Message-specific data
type ApiHandlerUnloadResponseJsonData struct {
	// Instance of the add-on package
	PackageName string `json:"packageName" yaml:"packageName"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerUnloadResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["packageName"]; !ok || v == nil {
		return fmt.Errorf("field packageName in ApiHandlerUnloadResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in ApiHandlerUnloadResponseJsonData: required")
	}
	type Plain ApiHandlerUnloadResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerUnloadResponseJsonData(plain)
	return nil
}

// Notice that an API handler is unloaded
type ApiHandlerUnloadResponseJson struct {
	// Message-specific data
	Data ApiHandlerUnloadResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and server to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerUnloadResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in ApiHandlerUnloadResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in ApiHandlerUnloadResponseJson: required")
	}
	type Plain ApiHandlerUnloadResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerUnloadResponseJson(plain)
	return nil
}
