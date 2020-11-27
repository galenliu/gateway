package addons

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
const ManifestFile = "manifest.yaml"
const ManifestFileJson = "manifest.json"

type AddonConfig struct {
	AddonManifest
	Enabled bool `json:"enabled"`
}

func NewAddonConfig(manifest AddonManifest) *AddonConfig {
	return &AddonConfig{AddonManifest: manifest, Enabled: true}
}

type AddonManifest struct {
	ID                      string               `json:"id"`
	Name                    string               `json:"name"`
	ShortName               string               `json:"short_name"`
	Author                  string               `json:"author"`
	Description             string               `json:"description,omitempty"`
	License                 string               `json:"license"`
	HomepageUrl             string               `json:"homepage_url,omitempty"`
	ManifestVersion         int                  `json:"manifest_version"`
	Version                 string               `json:"version"`
	Options                 Option               `json:"options,omitempty"`
	GatewaySpecificSettings map[string]WebThings `json:"gateway_specific_settings"`
	Enable                  bool                 `json:"_"`
}

type Option struct {
	Default map[string]interface{} `json:"default"`
	Schema  Schema                 `json:"schema"`
}
type Schema struct {
	Type       string              `json:"type"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}
type Property struct {
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Default string   `json:"default,omitempty"`
	Items   Schema   `json:"items,omitempty"`
	Enum    []string `json:"enum,omitempty"`
}

type WebThings struct {
	Exec             string `json:"exec"`
	PrimaryType      string `json:"primary_type"`
	StrictMaxVersion string `json:"strict_max_version"`
	StrictMinVersion string `json:"strict_min_version"`
}

func LoadManifest(addonDir string, packetId string) (*AddonManifest, error) {

	destPath := path.Join(addonDir, packetId)

	//load manifest.json
	manifest, err := loadManifestJson(destPath)
	if err != nil || &manifest == nil {
		e := fmt.Errorf("connot load manifest.js form %s", destPath)
		return nil, e
	}

	//First verify manifest version
	if manifest.ManifestVersion != ManifestVersion {
		err = fmt.Errorf("the manifest version(%v) for addon :%v does not match version",
			manifest.ManifestVersion, ManifestVersion)
		return nil, err
	}

	//verify that id in packet matches packetId
	if manifest.ID != packetId {
		err = fmt.Errorf("ID:%s from the manfest file,doesn't match ID for packetId: %s ",
			manifest.ID, packetId)
		return nil, err
	}

	//TODO :checksum every file.
	//TODO: Verify that manifest filed schema
	return manifest, err
}

func loadYaml(dirName string, in *AddonManifest) error {
	f, err := ioutil.ReadFile(path.Join(dirName, ManifestFile))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, in)
	return err
}

func loadManifestJson(dirName string) (addonManifest *AddonManifest, err error) {
	f, err := ioutil.ReadFile(path.Join(dirName, ManifestFileJson))

	var manifest AddonManifest
	err = jsoniter.Unmarshal(f, &manifest)
	return &manifest, err
}
