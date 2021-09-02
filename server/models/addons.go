package models

import (
	"github.com/galenliu/gateway/pkg/logging"
)

type AddonStore interface {
	LoadAddonSetting(id string) (value string, err error)
	StoreAddonSetting(id, value string) error
	LoadAddonConfig(id string) (value string, err error)
	StoreAddonsConfig(id, value string) error
}

type AddonManager interface {
	GetInstallAddons() []byte
	EnableAddon(addonId string) error
	DisableAddon(addonId string) error
	InstallAddonFromUrl(id, url, checksum string) error
	UnloadAddon(id string) error
	LoadAddon(id string) error
	UninstallAddon(id string, disabled bool) error
	GetAddonLicense(addonId string) (string, error)
	AddonEnabled(addonId string) bool
}

type AddonsModel struct {
	Store  AddonStore
	logger logging.Logger
	AddonManager
}

func NewAddonsModel(m AddonManager, store AddonStore, log logging.Logger) *AddonsModel {
	a := &AddonsModel{}
	a.logger = log
	a.Store = store
	a.AddonManager = m
	return a
}
