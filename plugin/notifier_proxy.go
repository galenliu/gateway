package plugin

import (
	"sync"
)

type Notifier struct {
	ID      string
	outlets sync.Map
}

func (n *Notifier) unload() {

}

func (m *Notifier) getOutlet(outletId string) *Outlet {
	a, ok := m.outlets.Load(outletId)
	outlet, ok := a.(*Outlet)
	if !ok {
		return nil
	}
	return outlet
}

func (m *Notifier) getOutlets() (outlets []*Outlet) {
	m.outlets.Range(func(key, value interface{}) bool {
		outlet, ok := value.(*Outlet)
		if ok {
			outlets = append(outlets, outlet)
		}
		return true
	})
	return
}
