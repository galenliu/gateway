package messages

import "fmt"
import "encoding/json"

// IPC messages between WebThings Gateway and add-ons
type SchemaJson struct {
	// The top-level message which encapsulates all message types
	Message SchemaJsonMessage `json:"message" yaml:"message"`
}

// The top-level message which encapsulates all message types
type SchemaJsonMessage map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *SchemaJson) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["message"]; !ok || v == nil {
		return fmt.Errorf("field message in SchemaJson: required")
	}
	type Plain SchemaJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SchemaJson(plain)
	return nil
}
