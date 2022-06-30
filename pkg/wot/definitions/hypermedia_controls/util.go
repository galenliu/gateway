package hypermedia_controls

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func JSONGetString(data []byte, field string, defaultValue string) string {
	if s := gjson.GetBytes(data, field); s.Type == gjson.String {
		return s.String()
	}
	return defaultValue
}

func JSONGetBool(data []byte, field string, defaultValue bool) bool {
	if b := gjson.GetBytes(data, field); b.Exists() {
		return b.Bool()
	}
	return defaultValue
}

func JSONGetFloat64(data []byte, field string, defaultValue float64) float64 {
	if f := gjson.GetBytes(data, field); f.Exists() {
		return f.Float()
	}
	return defaultValue
}

func JSONGetUint64(data []byte, field string, def uint64) uint64 {
	if f := gjson.GetBytes(data, field); f.Exists() {
		return f.Uint()
	}
	return def
}

func JSONGetMap(data []byte, field string) map[string]string {
	var m map[string]string
	if result := gjson.GetBytes(data, field); result.Exists() {
		err := json.Unmarshal([]byte(result.Raw), &m)
		if err == nil && m != nil && len(m) > 0 {
			return m
		}
	}
	return nil
}

func JSONGetInteger(data []byte, sep string) *Integer {
	if v := gjson.GetBytes(data, sep); v.Exists() {
		i := Integer(v.Int())
		return &i
	}
	return nil
}

func JSONGetArray(data []byte, field string) []string {
	var arr []string
	if result := gjson.GetBytes(data, field); result.Exists() {
		list := result.Array()
		for _, l := range list {
			arr = append(arr, l.String())
		}
		return arr
	}
	return nil
}
