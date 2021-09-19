package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type ThingDescription struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	AtContext   string `json:"@context"`
	AtType      string `json:"@type"`
	Description string `json:"description"`
	Base        string `json:"base"`
	BaseHref    string `json:"baseHref"`
	Href        string `json:"href"`

	Properties map[string]*Property `json:"properties"`
	Actions    map[string]*Action   `json:"actions"`
	Events     map[string]*Event    `json:"events"`

	Links               []Link             `json:"links,omitempty"`
	CredentialsRequired bool               `json:"credentialsRequired"`
	FloorplanVisibility bool               `json:"floorplanVisibility"`
	FloorplanX          uint               `json:"floorplanX"`
	FloorplanY          uint               `json:"floorplanY"`
	LayoutIndex         uint               `json:"layoutIndex"`
	SelectedCapability  string             `json:"selectedCapability"`
	IconHref            string             `json:"iconHref"`
	IconData            IconData           `json:"iconData"`
	Security            string             `json:"security"`
	SecurityDefinitions SecurityDefinition `json:"securityDefinitions"`
	GroupId             string             `json:"group_id"`
}

type IconData struct {
	Data string `json:"data"`
	Mime string `json:"mime"`
}

type OAuth2 struct {
	Scheme        string   `json:"scheme"`
	Flow          string   `json:"flow"`
	Authorization string   `json:"authorization"`
	Token         string   `json:"token"`
	Scopes        []string `json:"scopes"`
}
type SecurityDefinition struct {
	Oauth2Sc OAuth2 `json:"oauth2_sc"`
}

type Link struct {
	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

type Property struct {
	Name        string        `json:"name,omitempty"`
	AtType      string        `json:"@type,omitempty"`
	Title       string        `json:"title,omitempty"`
	Type        string        `json:"type"`
	Unit        string        `json:"unit,omitempty"`
	Description string        `json:"description,omitempty"`
	Minimum     interface{}   `json:"minimum,omitempty"`
	Maximum     interface{}   `json:"maximum,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	ReadOnly    bool          `json:"readOnly,omitempty"`
	MultipleOf  float64       `json:"multipleOf,omitempty"`
	Links       []Link        `json:"links,omitempty"`
	Value       interface{}   `json:"value,omitempty"`
}
type Action struct {
	AtType      string      `json:"@type,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Links       []Link      `json:"links,omitempty"`
	Input       interface{} `json:"input"`
}

type Event struct {
	AtType      string        `json:"@type,omitempty"`
	Name        string        `json:"name,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Links       []Link        `json:"links"`
	Type        string        `json:"type"`
	Unit        string        `json:"unit"`
	Minimum     interface{}   `json:"minimum"`
	Maximum     interface{}   `json:"maximum"`
	MultipleOf  float64       `json:"multipleOf"`
	Enum        []interface{} `json:"enum"`
}

type PIN struct {
	Required bool        `json:"required,omitempty"`
	Pattern  interface{} `json:"pattern,omitempty"`
}

type Thing struct {
	*wot.Thing

	//The configuration  of the device
	Pin                 *PIN `json:"pin,omitempty"`
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconHref           string `json:"iconHref,omitempty"`

	FloorplanVisibility bool `json:"floorplanVisibility"`
	FloorplanX          uint `json:"floorplanX"`
	FloorplanY          uint `json:"floorplanY"`
	LayoutIndex         uint `json:"layoutIndex"`

	Security            string             `json:"security"`
	SecurityDefinitions SecurityDefinition `json:"securityDefinitions"`
	GroupId             string             `json:"group_id"`
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(id string, description string) (thing *Thing, err error) {
	if id == "" || description == "" {
		return nil, fmt.Errorf("id or description err")
	}
	data := []byte(description)
	var p PIN
	json.Get(data, "pin").ToVal(&p)

	t := Thing{
		Thing:               wot.NewThingFromString(description),
		Pin:                 &p,
		CredentialsRequired: json.Get(data, "credentialsRequired").ToBool(),
		SelectedCapability:  json.Get(data, "selectedCapability").ToString(),
		Connected:           json.Get(data, "connected").ToBool(),
		IconHref:            json.Get(data, "iconHref").ToString(),
		FloorplanVisibility: false,
		FloorplanX:          0,
		FloorplanY:          0,
		LayoutIndex:         0,
		Security:            "",
		SecurityDefinitions: SecurityDefinition{},
		GroupId:             "",
	}

	if len(t.AtType) < 1 || t.ID == "" {
		return nil, fmt.Errorf("@type or id err")
	}

	if t.SelectedCapability == "" {
		t.SelectedCapability = t.AtType[0]
	}

	if !util.In(t.SelectedCapability, t.AtType) {
		return nil, fmt.Errorf("selectedCapability err")
	}
	return &t, nil
}
