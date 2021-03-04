package plugin

import (
	"addon"
	"fmt"
	"gateway/pkg/database"
	"gateway/server/models/thing"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"path"
)

const ManifestVersion = 1
const FileName = "manifest.json"

func GetAddonKey(id string) string {
	return fmt.Sprintf("addons.%s", id)
}

type ManifestJson struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ShortName       string `json:"short_name,omitempty"`
	Author          string `json:"author"`
	Description     string `json:"description,omitempty"`
	License         string `json:"license"`
	HomepageUrl     string `json:"homepage_url,omitempty"`
	ManifestVersion int    `json:"manifest_version"`
	Version         string `json:"version"`
	Options         struct {
		Default interface{} `json:"default,omitempty"`
		Schema  *struct {
			Type       string      `json:"type,omitempty"`
			Required   []string    `json:"required,omitempty"`
			Properties interface{} `json:"properties,omitempty"`
		} `json:"schema,,omitempty"`
	} `json:"options,omitempty"`
	GatewaySpecificSettings struct {
		WebThings struct {
			Exec             string `json:"exec"`
			PrimaryType      string `json:"primary_type"`
			StrictMaxVersion string `json:"strict_max_version"`
			StrictMinVersion string `json:"strict_min_version"`
		} `json:"webthings"`
	} `json:"gateway_specific_settings"`
	Enable bool `json:"-"`
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
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	ShortName   string      `json:"short_name"`
	Author      string      `json:"author"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	HomepageUrl string      `json:"homepage_url"`
	Version     string      `json:"version"`
	Schema      interface{} `json:"schema,omitempty"`
	Exec        string      `json:"exec"`
	Enabled     bool        `json:"enabled"`
	PrimaryType string      `json:"primary_type"`
}

func (addonInfo *AddonInfo) UpdateFromDB() error {

	d, err := json.MarshalToString(addonInfo)
	if err != nil {
		return err
	}
	err = database.SetSetting(GetAddonKey(addonInfo.ID), d)
	return nil
}

func (addonInfo *AddonInfo) UpdateAddonInfoToDB(enable bool) error {
	addonInfo.Enabled = enable
	s, err := json.MarshalToString(addonInfo)
	if err != nil {
		return err
	}
	return database.SetSetting(GetAddonKey(addonInfo.ID), s)
}

func GetAddonInfoFromDB(id string) (*AddonInfo, error) {
	value, err := database.GetSetting(GetAddonKey(id))
	if err != nil {
		return nil, err
	}
	var a AddonInfo
	err = json.UnmarshalFromString(value, a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func loadManifest(destPath, packetId string) (*AddonInfo, error) {

	//load manifest.json\
	f, err := ioutil.ReadFile(path.Join(destPath, FileName))
	if err != nil {
		return nil, err
	}
	var manifest ManifestJson
	err = json.Unmarshal(f, &manifest)
	if err != nil {
		return nil, err
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

		Exec:        manifest.GatewaySpecificSettings.WebThings.Exec,
		Enabled:     true,
		PrimaryType: manifest.GatewaySpecificSettings.WebThings.PrimaryType,
	}
	if manifest.Options.Schema != nil {
		addonInfo.Schema = manifest.Options.Schema
	}
	return &addonInfo, nil
}

func asThing(device *addon.Device) *thing.Thing {

	t := thing.Thing{
		ID:                  device.ID,
		AtContext:           device.AtContext,
		AtType:              device.AtType,
		Title:               device.Title,
		Description:         device.Description,
		BaseHref:            fmt.Sprintf("/thing/%s", device.ID),
		Href:                "",
		CredentialsRequired: device.CredentialsRequired,
		Properties:          nil,
		Actions:             nil,
		Events:              nil,
		Connected:           false,
	}

	t.Properties = make(map[string]*thing.Property)

	for _, prop := range device.Properties {
		t.Properties[prop.Name] = &thing.Property{
			Name:        prop.Name,
			AtType:      prop.AtType,
			Type:        prop.Type,
			Title:       prop.Title,
			Description: prop.Description,
			Unit:        prop.Unit,
			ReadOnly:    false,
			Visible:     prop.Visible,
			Minimum:     prop.Minimum,
			Maximum:     prop.Maximum,
			Value:       prop.Value,
			Enum:        prop.Enum,
			Links:       nil,
			Href:        fmt.Sprintf("/thing/%s/properties/%s", device.ID, prop.Name),
			ThingId:     t.ID,
		}
	}
	return &t
}
