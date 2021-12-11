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
	var m map[string]string
	if result := json.Get(data, field); result.LastError() == nil {
		result.ToVal(&m)
		if m != nil && len(m) > 0 {
			return m
		}
	}
	return nil
}

func JSONGetInteger(data []byte, sep string) *Integer {
	if v := json.Get(data, sep); v.LastError() == nil {
		i := Integer(v.ToInt64())
		return &i
	}
	return nil
}

func JSONGetArray(data []byte, field string) []string {
	var arr []string
	if result := json.Get(data, field); result.LastError() == nil {
		result.ToVal(&arr)
		if &arr != nil && len(arr) > 0 {
			return arr
		}
		return nil
	}
	return nil
}
