// Package plugin
// @Description:
package plugin

import (
	"fmt"
	json "github.com/json-iterator/go"
	"path"
)

const ManifestVersion = 1
const manifestFileName = "manifest.json"

type AddonsStore interface {
	LoadAddonSetting(key string) (string, error)
	StoreAddonSetting(key, value string) error
	LoadAddonConfig(key string) (string, error)
	StoreAddonsConfig(key string, value any) error
}

type Addon struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ShortName               string `json:"short_name"`
	Author                  string `json:"author"`
	Description             string `json:"description，omitempty"`
	License                 string `json:"license,omitempty"`
	HomepageUrl             string `json:"homepage_url"`
	Version                 string `json:"version"`
	Schema                  Schema `json:"schema,omitempty"`
	Exec                    string `json:"exec"`
	Enabled                 bool   `json:"enabled"`
	PrimaryType             string `json:"primary_type"`
	ContentScripts          string `json:"content_scripts,omitempty"`
	WSebAccessibleResources string `json:"web_accessible_resources,omitempty"`
	store                   AddonsStore
}

func NewAddonSettingFromString(str string, store AddonsStore) *Addon {
	var a Addon
	err := json.UnmarshalFromString(str, &a)
	if err != nil {
		return nil
	}
	a.store = store
	return &a
}

func NewAddonSettingFromManifest(manifest *ManifestJson, store AddonsStore) *Addon {
	addonInfo := &Addon{
		ID:                      manifest.ID,
		Name:                    manifest.Name,
		ShortName:               manifest.ShortName,
		Author:                  manifest.Author,
		Description:             manifest.Description,
		License:                 manifest.License,
		HomepageUrl:             manifest.HomepageUrl,
		Version:                 manifest.Version,
		ContentScripts:          manifest.ContentScripts,
		WSebAccessibleResources: manifest.WSebAccessibleResources,
		Exec:                    manifest.GatewaySpecificSettings.WebThings.Exec,
		Enabled:                 true,
		Schema:                  manifest.Options.Schema,
		PrimaryType:             manifest.GatewaySpecificSettings.WebThings.PrimaryType,
	}
	addonInfo.store = store
	return addonInfo
}

func (a *Addon) SetEnabled(disabled bool) error {
	if a.Enabled == disabled {
		return nil
	}
	a.Enabled = disabled
	err := a.save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Addon) save() error {
	str, err := json.MarshalToString(a)
	if err != nil {
		return err
	}
	return a.store.StoreAddonSetting(a.ID, str)
}

// LoadManifest
//  @Description:
//  @param destPath
//  @param packetId
//  @param store
//  @return *Addon
//  @return interface{}:addon default config
//  @return error
func LoadManifest(destPath, packetId string, store AddonsStore) (*Addon, any, error) {

	//load manifest.json
	manifest, err := readManifest(path.Join(destPath, manifestFileName))
	if err != nil {
		return nil, nil, err
	}
	//First verify manifest version
	if manifest.ManifestVersion != ManifestVersion {
		err = fmt.Errorf("the manifest version(%v) for addon :%v does not match version",
			manifest.ManifestVersion, ManifestVersion)
		return nil, nil, err
	}

	//verify that Id in packet matches packetId
	if manifest.ID != packetId {
		err = fmt.Errorf("Id:%s from the manfest file,doesn't match Id for packetId: %s ",
			manifest.ID, packetId)
		return nil, nil, err
	}

	//var min = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion
	//var max = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion

	//TODO :checksum every file.
	//TODO: Verify that manifest filed schema
	addonInfo := NewAddonSettingFromManifest(manifest, store)

	if !manifest.GatewaySpecificSettings.WebThings.Enable {
		addonInfo.Enabled = true
	}
	return addonInfo, manifest.Options.Default, nil
}
