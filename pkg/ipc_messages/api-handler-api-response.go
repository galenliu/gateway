// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package messages

import "fmt"
import "encoding/json"

// Outgoing response
type ApiHandlerApiResponseJsonDataResponse struct {
	// Body content
	Content any `json:"content,omitempty" yaml:"content,omitempty"`

	// Content-Type of the response body
	ContentType any `json:"contentType,omitempty" yaml:"contentType,omitempty"`

	// HTTP status code
	Status int `json:"status" yaml:"status"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerApiResponseJsonDataResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["status"]; !ok || v == nil {
		return fmt.Errorf("field status in ApiHandlerApiResponseJsonDataResponse: required")
	}
	type Plain ApiHandlerApiResponseJsonDataResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerApiResponseJsonDataResponse(plain)
	return nil
}

// Message-specific data
type ApiHandlerApiResponseJsonData struct {
	// ID of the request message
	MessageId int `json:"messageId" yaml:"messageId"`

	// Instance of the add-on package
	PackageName string `json:"packageName" yaml:"packageName"`

	// ID of the plugin
	PluginId string `json:"pluginId" yaml:"pluginId"`

	// Outgoing response
	Response ApiHandlerApiResponseJsonDataResponse `json:"response" yaml:"response"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerApiResponseJsonData) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["messageId"]; !ok || v == nil {
		return fmt.Errorf("field messageId in ApiHandlerApiResponseJsonData: required")
	}
	if v, ok := raw["packageName"]; !ok || v == nil {
		return fmt.Errorf("field packageName in ApiHandlerApiResponseJsonData: required")
	}
	if v, ok := raw["pluginId"]; !ok || v == nil {
		return fmt.Errorf("field pluginId in ApiHandlerApiResponseJsonData: required")
	}
	if v, ok := raw["response"]; !ok || v == nil {
		return fmt.Errorf("field response in ApiHandlerApiResponseJsonData: required")
	}
	type Plain ApiHandlerApiResponseJsonData
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerApiResponseJsonData(plain)
	return nil
}

// Response message from an API handler request
type ApiHandlerApiResponseJson struct {
	// Message-specific data
	Data ApiHandlerApiResponseJsonData `json:"data" yaml:"data"`

	// The message type, used by the IPC client and api to differentiate messages
	MessageType int `json:"messageType" yaml:"messageType"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ApiHandlerApiResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["data"]; !ok || v == nil {
		return fmt.Errorf("field data in ApiHandlerApiResponseJson: required")
	}
	if v, ok := raw["messageType"]; !ok || v == nil {
		return fmt.Errorf("field messageType in ApiHandlerApiResponseJson: required")
	}
	type Plain ApiHandlerApiResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ApiHandlerApiResponseJson(plain)
	return nil
}
