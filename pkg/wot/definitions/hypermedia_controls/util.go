package hypermedia_controls

import json "github.com/json-iterator/go"

func JSONGetString(data []byte, field string, defaultValue string) string {
	if s := json.Get(data, field); s.ValueType() == json.StringValue {
		return s.ToString()
	}
	return defaultValue
}

func JSONGetBool(data []byte, field string, defaultValue bool) bool {
	if b := json.Get(data, field); b.ValueType() == json.BoolValue {
		return b.ToBool()
	}
	return defaultValue
}

func JSONGetFloat64(data []byte, field string, defaultValue float64) float64 {
	if f := json.Get(data, field); f.ValueType() == json.NumberValue {
		return f.ToFloat64()
	}
	return defaultValue
}

func JSONGetUint64(data []byte, field string, def uint64) uint64 {
	if f := json.Get(data, field); f.ValueType() == json.NumberValue {
		return f.ToUint64()
	}
	return def
}

func JSONGetMap(data []byte, field string) map[string]string {
	var v map[string]string
	json.Get(data, field).ToVal(&v)
	if len(v) == 0 {
		return nil
	}
	return v
}

func JSONGetArray(data []byte, field string) []string {
	if result := json.Get(data, field); result.ValueType() == json.ArrayValue {
		return result.Keys()
	}

	if result := json.Get(data, field); result.ValueType() == json.StringValue {
		return []string{result.ToString()}
	}
	return nil
}
