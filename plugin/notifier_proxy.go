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

func (m *Notifier) getOutlet(outletId string) *models.Outlet {
	a, ok := m.outlets.Load(outletId)
	outlet, ok := a.(*models.Outlet)
	if !ok {
		return nil
	}
	return outlet
}

func (m *Notifier) getOutlets() (outlets []*models.Outlet) {
	m.outlets.Range(func(key, value interface{}) bool {
		outlet, ok := value.(*models.Outlet)
		if ok {
			outlets = append(outlets, outlet)
		}
		return true
	})
	return
}
