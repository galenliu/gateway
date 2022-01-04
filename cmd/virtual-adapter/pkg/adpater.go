package pkg

import (
	"github.com/galenliu/gateway/pkg/addon"
	"time"
)

type VirtualAdapter struct {
	*addon.Adapter
}

func NewVirtualAdapter(manager *addon.Manager, adapterId, name string) *VirtualAdapter {
	v := &VirtualAdapter{
		addon.NewAdapter(manager, adapterId, name),
	}
	return v
}

func (a VirtualAdapter) StartPairing(timeout time.Duration) {

}
