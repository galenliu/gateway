package models

import (
	"github.com/galenliu/gateway/pkg/logging"
)

type AddonStore interface {
	LoadAddonSetting(id string) (value string, err error)
	StoreAddonSetting(id, value string) error
	LoadAddonConfig(id string) (value string, err error)
	StoreAddonsConfig(id string, value interface{}) error
}

type AddonsModel struct {
	Store  AddonStore
	logger logging.Logger
}

func NewAddonsModel(store AddonStore, log logging.Logger) *AddonsModel {
	a := &AddonsModel{}
	a.logger = log
	a.Store = store
	return a
}
