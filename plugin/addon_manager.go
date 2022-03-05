package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/melbahja/got"
	"os"
	"path"
	"strings"
	"time"
)

// GetInstallAddons  获取已安装的add-on
func (m *Manager) GetInstallAddons() any {
	return m.getInstallAddons()
}

func (m *Manager) EnableAddon(packageId string) error {
	addonInfo := m.getInstallAddon(packageId)
	if addonInfo == nil {
		return fmt.Errorf("can't find addon %s", packageId)
	}
	err := addonInfo.SetEnabled(true)
	if err != nil {
		return err
	}
	err = m.LoadAddon(packageId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) DisableAddon(packageId string) error {
	addon := m.getInstallAddon(packageId)
	if addon == nil {
		return fmt.Errorf("can't find addon %s", packageId)
	}
	err := addon.SetEnabled(false)
	if err != nil {
		return err
	}
	err = m.UnloadAddon(packageId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) AddonEnabled(addonId string) bool {
	addon := m.getInstallAddon(addonId)
	return addon.Enabled
}

func (m *Manager) InstallAddonFromUrl(id, url, checksum string) error {

	list := strings.SplitAfter(url, "/")
	fileName := list[len(list)-1]

	file := path.Join(os.TempDir(), fileName)
	m.logger.Infof("fetching add-on %s as %s", url, file)
	g := got.New()
	err := g.Download(url, file)
	if err != nil {
		return fmt.Errorf("download %s error: %s", file, err.Error())
	}
	defer func() {
		err := os.Remove(file)
		if err != nil {
			m.logger.Infof("remove temp file failed ,err:%s", err.Error())
		}
	}()

	err = os.Chmod(file, 777)
	if err != nil {
		m.logger.Infof("chmod err,package id:%s,err:%v", id, err.Error())
		return err
	}
	if !util.CheckSum(file, checksum) {
		m.logger.Infof("checksum err,package id:%s", id)
		return fmt.Errorf("checksum err,package id:%s", id)
	}
	m.logger.Infof("download %s successful", id)
	err = m.installAddon(id, file)
	if err != nil {
		m.logger.Infof("install %s failed error: %s", id, err.Error())
		return err
	}
	return nil
}

func (m *Manager) UninstallAddon(packetId string, disable bool) error {
	// delete addon form manager
	defer func() {
		if disable {
			addon := m.getInstallAddon(packetId)
			if addon != nil {
				err := addon.DeleteSettingAndConfig()
				if err != nil {
					m.logger.Infof(err.Error())
				}
			}
		}
		m.installAddons.Delete(packetId)
	}()

	err := m.UnloadAddon(packetId)
	if err != nil {
		fmt.Printf(err.Error())
	}
	err = util.RemoveDir(path.Join(m.config.AddonsDir, packetId))
	if err != nil {
		m.logger.Errorf("delete plugin addon dir failed err: %s", err.Error())
	}
	err = util.RemoveDir(path.Join(m.config.UserProfile.DataDir, packetId))
	if err != nil {
		m.logger.Errorf("delete plugin data dir failed err:%s", err.Error())
	}
	return nil
}

func (m *Manager) LoadAddon(packageId string) error {
	m.logger.Infof("starting loading addon %s.", packageId)
	var addonInfo *Addon
	var obj any
	var err error
	packageDir := path.Join(m.config.AddonsDir, packageId)
	addonInfo, obj, err = LoadManifest(packageDir, packageId, m.storage)
	if err != nil {
		return fmt.Errorf("load file %s  manifest.json  err:", packageId)
	}
	saved, err := m.storage.LoadAddonSetting(packageId)
	if err == nil && saved != "" {
		addonInfo = NewAddonSettingFromString(saved, m.storage)
	} else {
		err = addonInfo.save()
		if err != nil {
			m.logger.Errorf("addon save err: %s", err.Error())
		}
	}
	addonConf, err := m.storage.LoadAddonConfig(packageId)
	if err != nil && addonConf == "" {
		if obj != nil {
			err := m.storage.StoreAddonsConfig(packageId, obj)
			if err != nil {
				m.logger.Errorf("store addon config err: %s", err.Error())
			}
		}
	}
	m.installAddons.Store(packageId, addonInfo)

	if addonInfo.ContentScripts != nil && addonInfo.WSebAccessibleResources != "" {
		var ext = Extension{
			Extensions: addonInfo.ContentScripts,
			Resources:  addonInfo.WSebAccessibleResources,
		}
		m.extensions.Store(addonInfo.ID, ext)
	}
	util.EnsureDir(m.logger, path.Join(m.config.UserProfile.DataDir, packageId))
	if addonInfo.Exec == "" {
		return fmt.Errorf("addon %s has not exec", addonInfo.ID)

	}
	if !addonInfo.Enabled {
		return fmt.Errorf("addon %s disabled", packageId)
	}
	m.pluginServer.loadPlugin(addonInfo.ID, packageDir, addonInfo.Exec)
	return nil
}

func (m *Manager) UnloadAddon(packageId string) error {
	if !m.addonsLoaded {
		return fmt.Errorf("the add-ons are not currently loaded, no need to notifyUnload")
	}
	plugin := m.getPlugin(packageId)
	if plugin == nil {
		return fmt.Errorf("The add-ons are not  register.")
	}
	plugin.unloadComponents()
	go time.AfterFunc(time.Duration(constant.UnloadPluginKillDelay)*time.Millisecond, func() {
		plugin.shutdown()
		//_, _ = m.extensions.LoadAndDelete(packageId)
	})
	return nil
}

func (m *Manager) GetAddonLicense(addonId string) (string, error) {
	addonDir := path.Join(m.config.AddonsDir, addonId)
	m.pluginServer.getPlugin(addonId)
	if addonDir == "" {
		return "", fmt.Errorf("addon not installed")
	}
	data, err := os.ReadFile(path.Join(addonDir, "LICENSE"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
