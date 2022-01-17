package virtual

import (
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

type Adapter struct {
	*proxy.Adapter
}

func NewVirtualAdapter(adapterId string) *Adapter {
	v := &Adapter{
		proxy.NewAdapter(adapterId, "Virtual"),
	}
	return v
}
