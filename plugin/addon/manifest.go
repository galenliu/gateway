package addon

import (
	json "github.com/json-iterator/go"
	"io/ioutil"
)

type Schema struct {
	Type       string      `json:"type,omitempty"`
	Required   []string    `json:"required,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

type ManifestJson struct {
	ID                      string `json:"ID"`
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
			Enable           bool   `json:"setEnabled"`
		} `json:"webthings"`
	} `json:"gateway_specific_settings"`
	Enable bool `json:"setEnabled"`
}

func readManifest(file string) (*ManifestJson, error) {

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var manifest ManifestJson
	err = json.Unmarshal(f, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, err
}
