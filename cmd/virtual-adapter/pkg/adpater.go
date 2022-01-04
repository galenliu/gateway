package pkg

import (
	"github.com/galenliu/gateway/pkg/addon"
	"time"
)

type VirtualAdapter struct {
	*addon.Adapter
}

func NewVirtualAdapter(adapterId, name string) *VirtualAdapter {
	v := &VirtualAdapter{
		addon.NewAdapter(adapterId, name),
	}
	return v
}

func (a VirtualAdapter) StartPairing(timeout time.Duration) {

}
