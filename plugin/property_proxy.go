package plugin

import (
	ipc "github.com/galenliu/smartassistant-ipc"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type PropertyProxy struct {
	*ipc.Property
}

func (proxy *PropertyProxy) gitName() string {
	return proxy.Name
}

func (proxy *PropertyProxy) AsDict() (d string) {
	d, e := json.MarshalToString(proxy)
	if e != nil {
		log.Warn("property marshal err", zap.Error(e))
	}
	return d
}
