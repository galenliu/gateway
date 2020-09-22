package plugin

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
const ManifestFile = "manifest.toml"

type Manifest struct {
	Author          string                 `toml:"author"`
	Description     string                 `toml:"description,omitempty"`
	Name            string                 `toml:"name"`
	ID              string                 `toml:"id"`
	ManifestVersion int                    `toml:"manifest_version"`
	HomepageUrl     string                 `toml:"homepage_url,omitempty"`
	Enabled         bool                   `toml:"enabled"`
	Options         map[string]interface{} `toml:"options"`
	Settings        Settings               `toml:"settings"`
}

type Settings struct {
	Exec       string `toml:"exec"`
	MinVersion string `toml:"min_version"`
	MaxVersion string `toml:"max_version"`
}

func LoadManifest(addonDir string, addonDirName string) (*Manifest, error) {
	var manifest Manifest
	packetDir := path.Join(addonDir, addonDirName)
	err := loadToml(packetDir, &manifest)
	if err != nil || &manifest == nil {
		return nil, err
	}

	//First verify manifest version
	if manifest.ManifestVersion != ManifestVersion {
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

func loadToml(dirName string, in *Manifest) error {
	f, err := ioutil.ReadFile(path.Join(dirName, ManifestFile))
	if err != nil {
		return err
	}
	err = toml.Unmarshal(f, in)
	return err
}
