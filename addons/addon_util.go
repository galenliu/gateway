package addons

import (
	"errors"
	"fmt"
	"gateway/pkg/database"
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
	"gorm.io/gorm"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
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
	ShortName               string               `json:"short_name,omitempty"`
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
	Schema  Schema                 `json:"schema,omitempty"`
}
type Schema struct {
	Type       string                `json:"type,omitempty"`
	Required   []string              `json:"required,omitempty"`
	Properties map[string]Properties `json:"properties,omitempty"`
}
type Properties struct {
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

//author: "bewee"
//description: "Tuya Smart Life IoT devices support"
//enabled: true
//exec: "{nodeLoader} {path}"
//homepage_url: "https://github.com/bewee/tuya-adapter"
//id: "tuya-adapter"
//name: "Tuya Smart Life"
//primary_type: "adapter"
//schema: {type: "object", required: ["devices", "timeout", "log"], properties: {,…}}
//version: "0.2.4"

type AddonInfo struct {
	gorm.Model
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	ShortName   string      `json:"short_name"`
	Author      string      `json:"author"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	HomepageUrl string      `json:"homepage_url"`
	Version     string      `json:"version"`
	Schema      interface{} `json:"schema,omitempty" gorm:"-"`
	Exec        string      `json:"exec"`
	Enabled     bool        `json:"enabled"`
	PrimaryType string      `json:"primary_type" gorm:"default: adapter"`
}

func (addonInfo *AddonInfo) UpdateOrCreateFormDb() error {

	db, err := database.GetDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&AddonInfo{})
	if err != nil {
		return err
	}
	var a =AddonInfo{}
	rst := db.Where("id = ?",addonInfo.ID).First(&a)
	if errors.Is(rst.Error,gorm.ErrRecordNotFound){
		db.Create(addonInfo)
		log.Info("create new addonInfo(%s) record", addonInfo.ID)
		return nil
	}
	addonInfo.Enabled = a.Enabled
	return nil
}

func (addonInfo *AddonInfo) UpdateAddonInfoToDB(enable bool) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}
	_ = db.AutoMigrate(&AddonInfo{})
	addonInfo.Enabled = enable
	tx := db.Save(addonInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetAddonsInfoFromDB() []AddonInfo {
	db, _ := database.GetDB()
	_ = db.AutoMigrate(AddonInfo{})
	var addons []AddonInfo
	db.Find(&addons)
	return addons
}
func GetAddonInfoByIDFromDB(packageId string) (*AddonInfo, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(AddonInfo{})
	if err != nil {
		return nil, err
	}
	var addonInfo AddonInfo
	tx := db.Find(&addonInfo, packageId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &addonInfo, nil
}

func LoadManifest(addonDir string, packetId string) (*AddonInfo, error) {

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

	addonInfo := AddonInfo{
		ID:          manifest.ID,
		Name:        manifest.Name,
		ShortName:   manifest.ShortName,
		Author:      manifest.Author,
		Description: manifest.Description,
		License:     manifest.License,
		HomepageUrl: manifest.HomepageUrl,
		Version:     manifest.Version,
		Schema:      manifest.Options.Schema,
		Exec:        manifest.GatewaySpecificSettings["webthings"].Exec,
		Enabled:     true,
		PrimaryType: manifest.GatewaySpecificSettings["webthings"].PrimaryType,
	}
	return &addonInfo, nil
}

func loadManifestJson(dirName string) (addonManifest *AddonManifest, err error) {
	f, err := ioutil.ReadFile(path.Join(dirName, ManifestFileJson))

	var manifest AddonManifest
	err = json.Unmarshal(f, &manifest)
	return &manifest, err
}
