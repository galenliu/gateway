package plugin

import (
	"fmt"
	json "github.com/json-iterator/go"
	"path"
)

const ManifestVersion = 1
const FileName = "manifest.json"

type AddonInfo struct {
	ID                      string `json:"ID"`
	Name                    string `json:"name"`
	ShortName               string `json:"short_name"`
	Author                  string `json:"author"`
	Description             string `json:"description"`
	License                 string `json:"license"`
	HomepageUrl             string `json:"homepage_url"`
	Version                 string `json:"version"`
	Schema                  Schema `json:"schema,omitempty"`
	Exec                    string `json:"exec"`
	Enabled                 bool   `json:"enabled"`
	PrimaryType             string `json:"primary_type"`
	ContentScripts          string `json:"content_scripts"`
	WSebAccessibleResources string `json:"web_accessible_resources"`
	store                   Store
	dir                     string
}

func NewAddonInfoFromManifest(manifest *ManifestJson, store Store) *AddonInfo {
	addonInfo := &AddonInfo{
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

func (a *AddonInfo) setEnabled(disabled bool) error {
	if a.Enabled == disabled {
		return nil
	}
	a.Enabled = disabled
	b, err := json.Marshal(a)
	err = a.store.SetAddonsSetting(a.ID, string(b))
	if err != nil {
		return err
	}
	return nil
}

func LoadManifest(destPath, packetId string, store Store) (*AddonInfo, interface{}, error) {

	//load manifest.json\
	manifest, err := ReadManifestJson(path.Join(destPath, FileName))
	if err != nil {
		return nil, nil, err
	}
	//First verify manifest version
	if manifest.ManifestVersion != ManifestVersion {
		err = fmt.Errorf("the manifest version(%v) for addon :%v does not match version",
			manifest.ManifestVersion, ManifestVersion)
		return nil, nil, err
	}

	//verify that ID in packet matches packetId
	if manifest.ID != packetId {
		err = fmt.Errorf("Id:%s from the manfest file,doesn't match Id for packetId: %s ",
			manifest.ID, packetId)
		return nil, nil, err
	}

	//var min = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion
	//var max = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion

	//TODO :checksum every file.
	//TODO: Verify that manifest filed schema
	addonInfo := NewAddonInfoFromManifest(manifest, store)
	addonInfo.dir = destPath
	if !manifest.GatewaySpecificSettings.WebThings.Enable {
		addonInfo.Enabled = true
	}
	return addonInfo, manifest.Options.Default, nil
}
