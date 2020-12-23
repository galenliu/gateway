package models

import "gateway/addons"

type Thing struct {
	ID           string   `json:"id"`
	AtContext    string   `json:"@context omitempty"`
	AtType       []string `json:"@type"`
	Title        string   `json:"title,required"`
	Titles       []string `json:"titles,omitempty"`
	Description  string   `json:"description,omitempty"`
	Descriptions []string `json:"descriptions,omitempty"`
	Version      string   `json:"version,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
	Support      string   `json:"support,omitempty"`

	Links               []string    `json:"links"`
	BaseHref            string      `json:"base_href"`
	Href                string      `json:"href"`
	Pin                 interface{} `json:"pin"`
	CredentialsRequired bool        `json:"credentials_required"`
	SelectedCapability  string      `json:"selected_capability"`

	Properties map[string]*Property `json:"properties"`
	//Actions    interface{} `json:"actions"`
	//Events     interface{} `json:"events"`
}

type thing struct {
	Thing
	connected bool
}

func NewThing(t Thing) *thing {
	th := &thing{Thing: t}
	return th
}

func deviceToThing(devices map[string]*addons.DeviceProxy) map[string]*Thing {

	var thingsMap = make(map[string]*Thing)
	for key, dev := range devices {
		thing := &Thing{
			ID:                  dev.ID,
			AtContext:           dev.AtContext,
			AtType:              dev.AtType,
			Title:               dev.Title,
			Titles:              dev.Titles,
			Description:         dev.Description,
			Descriptions:        dev.Descriptions,
			Version:             dev.Version,
			Created:             dev.Created,
			Modified:            dev.Modified,
			Support:             dev.Support,
			Links:               nil,
			BaseHref:            "",
			Href:                "",
			Pin:                 nil,
			CredentialsRequired: false,
			SelectedCapability:  "",
			Properties:          nil,
		}
		var props = make(map[string]*Property)
		for _, p := range dev.Properties {

			thingProp := devPropToThingProp(p)
			props[thingProp.Name] = thingProp

		}
		thing.Properties = props
		thingsMap[key] = thing
	}

	return thingsMap
}
