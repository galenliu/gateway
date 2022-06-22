package models

type AddonStore interface {
	LoadAddonSetting(id string) (value string, err error)
	StoreAddonSetting(id, value string) error
	LoadAddonConfig(id string) (value string, err error)
	StoreAddonsConfig(id string, value any) error
}

type AddonsModel struct {
	Store AddonStore
}

func NewAddonsModel(store AddonStore) *AddonsModel {
	a := &AddonsModel{}

	a.Store = store
	return a
}
