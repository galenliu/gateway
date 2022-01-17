package virtual

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"time"
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

func (v *Adapter) StartPairing(t <-chan time.Time) {
	fmt.Printf("adapter: %s StartPairing", v.GetId())
}
