package plugin

import (
	"addon"
	"fmt"
	"gateway/pkg/database"
	"gateway/pkg/util"
	"gateway/server/models"
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

func (addonInfo *AddonInfo) UpdateFromDB() error {
	savedAddonInfo := GetAddonInfoFromDB(addonInfo.ID)
	if savedAddonInfo != nil {
		addonInfo.Enabled = savedAddonInfo.Enabled
	}

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

func GetAddonInfoFromDB(id string) *AddonInfo {
	value, err := database.GetSetting(GetAddonKey(id))
	if err != nil {
		return nil
	}
	var a AddonInfo
	err = json.UnmarshalFromString(value, a)
	if err != nil {
		return nil
	}
	return &a
}

func loadManifest(destPath, packetId string) (*AddonInfo, *interface{}, error) {

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

func MarshalWebThing(data []byte) (*models.Thing, error) {

	id := json.Get(data, "id").ToString()
	if id == "" {
		return nil, fmt.Errorf("id necessary")
	}
	title := json.Get(data, "title").ToString()
	if title == "" {
		title = id
	}
	id = fmt.Sprintf("/things/%s", id)

	var atContext []string
	json.Get(data, "@context").ToVal(&atContext)

	var atType []string
	json.Get(data, "@type").ToVal(&atType)

	t := &models.Thing{
		AtContext:           atContext,
		Title:               title,
		ID:                  id,
		AtType:              atType,
		Description:         json.Get(data, "description").ToString(),
		Properties:          nil,
		Actions:             nil,
		Events:              nil,
		Forms:               nil,
		CredentialsRequired: json.Get(data, "credentialsRequired").ToBool(),
	}

	var pin *addon.PIN
	json.Get(data, "pin").ToVal(&pin)
	if pin != nil {
		t.Pin = *pin
	}
	var props map[string]addon.Property
	json.Get(data, "properties").ToVal(&props)
	if len(props) > 0 {
		t.Properties = make(map[string]*models.Property)
		for n, p := range props {
			prop := &models.Property{
				Name:        n,
				AtType:      p.AtType,
				Type:        p.Type,
				Title:       p.Title,
				Description: p.Description,
				Unit:        p.Unit,
				ReadOnly:    p.ReadOnly,
				Visible:     p.Visible,
				Minimum:     p.Minimum,
				Maximum:     p.Maximum,
				Enum:        p.Enum,
				ThingId:     t.ID,
			}
			prop.Forms = append(prop.Forms, util.NewForm("href", fmt.Sprintf("%s/properties/%s", t.ID, prop.Name)))
			t.Properties[n] = prop
		}
		t.Forms = append(t.Forms, util.NewForm("rel", "alternate", "mediaType", "text/html", "href", fmt.Sprintf("/things/%s", id)))
		t.Forms = append(t.Forms, util.NewForm("rel", "alternate", "href", fmt.Sprintf("/things/%s/", id)))
	}
	return t, nil
}

func UnmarshalDevice(d []byte) (*addon.Device, error) {

	id := json.Get(d, "id").ToString()
	if id == "" {
		return nil, fmt.Errorf("device id lost")
	}
	title := json.Get(d, "title").ToString()
	if title == "" {
		title = id
	}

	atContext := json.Get(d, "@context").Keys()
	if len(atContext) == 0 {
		t := json.Get(d, "@context").ToString()
		if t != "" {
			atContext = append(atContext, t)
		}
	}

	var atType []string
	json.Get(d, `@type`).ToVal(&atType)
	if len(atType) == 0 {
		return nil, fmt.Errorf("@type lost")
	}

	var properties map[string]*addon.Property
	json.Get(d, "properties").ToVal(&properties)

	var actions map[string]*addon.Action
	json.Get(d, "actions").ToVal(&actions)

	var events map[string]*addon.Event
	json.Get(d, "actions").ToVal(&events)

	var pin *addon.PIN
	json.Get(d, "pin").ToVal(&pin)

	device := &addon.Device{
		ID:                  id,
		AtContext:           atContext,
		Title:               title,
		AtType:              atType,
		Description:         json.Get(d, "description").ToString(),
		CredentialsRequired: json.Get(d, "credentialsRequired").ToBool(),
		Pin:                 addon.PIN{},
		AdapterId:           json.Get(d, "adapterId").ToString(),
	}
	if len(properties) > 0 {
		device.Properties = make(map[string]addon.IProperty)
		for n, p := range properties {
			p.DeviceId = id
			device.Properties[n] = p
		}
	}

	if len(events) > 0 {
		device.Events = make(map[string]*addon.Event)
		for n, e := range events {
			e.DeviceId = id
			device.Events[n] = e
		}
	}

	if len(actions) > 0 {
		device.Actions = make(map[string]*addon.Action)
		for n, a := range actions {
			a.DeviceId = id
			device.Actions[n] = a
		}
	}

	if pin != nil {
		device.Pin = *pin
	}

	return device, nil
}
