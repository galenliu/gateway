package plugin

import (
	"github.com/galenliu/gateway/plugin/internal"
	"sync"
)

type Notifier struct {
	ID      string
	outlets sync.Map
}

func (n *Notifier) unload() {

}

func (m *Notifier) getOutlet(outletId string) *internal.Outlet {
	a, ok := m.outlets.Load(outletId)
	outlet, ok := a.(*internal.Outlet)
	if !ok {
		return nil
	}
	return outlet
}

func (m *Notifier) getOutlets() (outlets []*internal.Outlet) {
	m.outlets.Range(func(key, value interface{}) bool {
		outlet, ok := value.(*internal.Outlet)
		if ok {
			outlets = append(outlets, outlet)
		}
		return true
	})
	return
}
