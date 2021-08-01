package plugin

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
const FileName = "manifest.json"

func GetAddonKey(id string) string {
	return fmt.Sprintf("addons.%s", id)
}

type Schema struct {
	Type       string      `json:"type,omitempty"`
	Required   []string    `json:"required,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

type ManifestJson struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ShortName               string `json:"short_name,omitempty"`
	Author                  string `json:"author"`
	Description             string `json:"description,omitempty"`
	License                 string `json:"license"`
	HomepageUrl             string `json:"homepage_url,omitempty"`
	ManifestVersion         int    `json:"manifest_version"`
	Version                 string `json:"version"`
	ContentScripts          string `json:"content_Scripts"`
	WSebAccessibleResources string `json:"web_accessible_resources"`
	Options                 struct {
		Default interface{} `json:"default"`
		Schema  Schema      `json:"schema"`
	}
	GatewaySpecificSettings struct {
		WebThings struct {
			Exec             string `json:"exec"`
			PrimaryType      string `json:"primary_type"`
			StrictMaxVersion string `json:"strict_max_version"`
			StrictMinVersion string `json:"strict_min_version"`
			Enable           bool   `json:"enable"`
		} `json:"webthings"`
	} `json:"gateway_specific_settings"`
	Enable bool `json:"enable"`
}

type AddonInfo struct {
	ID                      string `json:"id"`
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
}

func NewAddonInfoFromString(des string) *AddonInfo {
	var addonInfo AddonInfo
	err := json.UnmarshalFromString(des, &addonInfo)
	if err != nil {
		return &addonInfo
	}
	return nil
}


func LoadManifest(destPath, packetId string) (*AddonInfo, *interface{}, error) {

	//load manifest.json\
	f, err := ioutil.ReadFile(path.Join(destPath, FileName))
	if err != nil {
		return nil, nil, err
	}
	var manifest ManifestJson
	err = json.Unmarshal(f, &manifest)
	if err != nil {
		return nil, nil, err
	}

	//First verify manifest version
	if manifest.ManifestVersion != ManifestVersion {
		err = fmt.Errorf("the manifest version(%v) for addon :%v does not match version",
			manifest.ManifestVersion, ManifestVersion)
		return nil, nil, err
	}

	//verify that id in packet matches packetId
	if manifest.ID != packetId {
		err = fmt.Errorf("Id:%s from the manfest file,doesn't match Id for packetId: %s ",
			manifest.ID, packetId)
		return nil, nil, err
	}

	//var min = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion
	//var max = manifest.GatewaySpecificSettings.WebThings.StrictMinVersion

	//TODO :checksum every file.
	//TODO: Verify that manifest filed schema
	addonInfo := AddonInfo{
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
	if !manifest.GatewaySpecificSettings.WebThings.Enable {
		addonInfo.Enabled = true
	}
	return &addonInfo, &manifest.Options.Default, nil
}
