package plugin

import (
	"github.com/galenliu/gateway/plugin/addon"
	"sync"
)

type Notifier struct {
	ID      string
	outlets sync.Map
}

func (n *Notifier) unload() {

}

func (m *Notifier) getOutlet(outletId string) *addon.Outlet {
	a, ok := m.outlets.Load(outletId)
	outlet, ok := a.(*addon.Outlet)
	if !ok {
		return nil
	}
	return outlet
}

func (m *Notifier) getOutlets() (outlets []*addon.Outlet) {
	m.outlets.Range(func(key, value interface{}) bool {
		outlet, ok := value.(*addon.Outlet)
		if ok {
			outlets = append(outlets, outlet)
		}
		return true
	})
	return
}
