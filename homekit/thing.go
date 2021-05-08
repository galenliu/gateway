package homekit

import (
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/homekit/things"
	json "github.com/json-iterator/go"
)

type Thing struct {
	id string
	*service.Service
}

func NewThing(data []byte) *Thing {
	var types []string
	json.Get(data, "@type").ToVal(&types)
	id := json.Get(data, "id").ToString()
	if types == nil {
		return nil
	}
	for _, typ := range types {
		switch typ {
		case Light:
			thing := things.NewLightBulb(data)
			thing.id = id
			return thing.Thing
		}
	}
	return nil
}
