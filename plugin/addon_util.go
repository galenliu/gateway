package plugin

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
const ManifestFile = "manifest.yaml"

type AddonConfig struct {
	AddonManifest
	Enabled bool `json:"enabled"`
}

func NewAddonConfig(manifest AddonManifest) *AddonConfig {
	return &AddonConfig{AddonManifest: manifest, Enabled: true}
}

type AddonManifest struct {
	ID               string      `yaml:"id"`
	Name             string      `yaml:"name"`
	Author           string      `yaml:"author"`
	Description      string      `yaml:"description,omitempty"`
	License          string      `yaml:"license"`
	HomepageUrl      string      `yaml:"homepage_url,omitempty"`
	ManifestVersion  int         `yaml:"manifest_version"`
	Version          int         `yaml:"version"`
	StrictMaxVersion int         `yaml:"strict_max_version"`
	StrictMinVersion int         `yaml:"strict_min_version"`
	Option           interface{} `yaml:"option"`
	Exec             string      `yaml:"exec"`
}

func LoadManifest(addonDir string, addonDirName string) (*AddonManifest, error) {
	var manifest AddonManifest
	packetDir := path.Join(addonDir, addonDirName)
	err := loadYaml(packetDir, &manifest)
	if err != nil || &manifest == nil {
		return nil, err
	}
	addonConfig := NewAddonConfig(manifest)

	//First verify manifest version
	if addonConfig.ManifestVersion != ManifestVersion {
		err = fmt.Errorf("the manifest version(%v) for addon :%v does not match version",
			manifest.ManifestVersion, ManifestVersion)
		log.Warn("", zap.Error(err))
		return nil, err
	}

	//verify that id in packet matches packetId
	if manifest.ID != addonDirName {
		err = fmt.Errorf("ID:%s from the manfest file,doesn't match ID for packetId: %s ",
			manifest.ID, addonDirName)
		return nil, err
	}
	return &manifest, err
}

func loadYaml(dirName string, in *AddonManifest) error {
	f, err := ioutil.ReadFile(path.Join(dirName, ManifestFile))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, in)
	return err
}
