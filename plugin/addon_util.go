package plugin

import (
	"addon"
	"fmt"
	"gateway/pkg/database"
	"gateway/pkg/util"
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

type Options struct {
	Default interface{} `json:"default,omitempty"`
	Schema  struct {
		Type       string      `json:"type,omitempty"`
		Required   []string    `json:"required,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
	} `json:"schema,,omitempty"`
}

type ManifestJson struct {
	ID                      string  `json:"id"`
	Name                    string  `json:"name"`
	ShortName               string  `json:"short_name,omitempty"`
	Author                  string  `json:"author"`
	Description             string  `json:"description,omitempty"`
	License                 string  `json:"license"`
	HomepageUrl             string  `json:"homepage_url,omitempty"`
	ManifestVersion         int     `json:"manifest_version"`
	Version                 string  `json:"version"`
	ContentScripts          string  `json:"content_Scripts"`
	WSebAccessibleResources string  `json:"web_accessible_resources"`
	Options                 Options `json:"options"`
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
	ID                      string      `json:"id"`
	Name                    string      `json:"name"`
	ShortName               string      `json:"short_name"`
	Author                  string      `json:"author"`
	Description             string      `json:"description"`
	License                 string      `json:"license"`
	HomepageUrl             string      `json:"homepage_url"`
	Version                 string      `json:"version"`
	Schema                  interface{} `json:"schema,omitempty"`
	Exec                    string      `json:"exec"`
	Enabled                 bool        `json:"enabled"`
	PrimaryType             string      `json:"primary_type"`
	ContentScripts          string      `json:"content_scripts"`
	WSebAccessibleResources string      `json:"web_accessible_resources"`
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

func loadManifest(destPath, packetId string) (*AddonInfo, *Options, error) {

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
		err = fmt.Errorf("ID:%s from the manfest file,doesn't match ID for packetId: %s ",
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
		Enabled:                 false,
		PrimaryType:             manifest.GatewaySpecificSettings.WebThings.PrimaryType,
	}
	if manifest.GatewaySpecificSettings.WebThings.Enable {
		addonInfo.Enabled = true
	}
	return &addonInfo, &manifest.Options, nil
}

func asWebThing(device *addon.Device) *thing.Thing {

	t := thing.Thing{
		ID:          fmt.Sprintf("/things/%s", device.ID),
		AtContext:   device.AtContext,
		AtType:      device.AtType,
		Title:       device.Title,
		Description: device.Description,
		Properties:  nil,
		Actions:     nil,
		Events:      nil,

		Connected:           false,
		CredentialsRequired: device.CredentialsRequired,
	}

	t.Properties = make(map[string]*thing.Property)
	if len(device.Properties) > 0 {

		f := util.NewForm("rel", "properties", "href", fmt.Sprintf("%s/properties", t.ID))
		t.Forms = append(t.Forms, f)

		for _, prop := range device.Properties {
			var thingProperty *thing.Property
			thingProperty = &thing.Property{
				Name:        prop.Name,
				AtType:      prop.AtType,
				Type:        prop.Type,
				Title:       prop.Title,
				Description: prop.Description,
				Unit:        prop.Unit,
				ReadOnly:    prop.ReadOnly,
				Visible:     prop.Visible,
				Minimum:     prop.Minimum,
				Maximum:     prop.Maximum,
				Value:       prop.Value,
				Enum:        prop.Enum,
				ThingId:     t.ID,
			}
			thingProperty.Forms = append(thingProperty.Forms, util.NewForm("href", fmt.Sprintf("%s/properties/%s", t.ID, prop.Name)))
			t.Properties[thingProperty.Name] = thingProperty
		}

	}
	t.Forms = append(t.Forms, util.NewForm("rel", "alternate", "mediaType", "text/html", "href", fmt.Sprintf("/things/%s", device.ID)))
	t.Forms = append(t.Forms, util.NewForm("rel", "alternate", "href", fmt.Sprintf("/things/%s/", device.ID)))
	return &t
}
