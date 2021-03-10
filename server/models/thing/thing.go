package thing

import (
	"fmt"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
)

type Thing struct {
	AtContext   []string `json:"@context,omitempty"`
	Title       string   `json:"title"`
	ID          string   `json:"id"`
	AtType      []string `json:"@type"`
	Description string   `json:"description,omitempty"`

	Properties map[string]*Property `json:"properties,omitempty"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`
	Forms      []util.Form          `json:"forms,omitempty"`

	//The configuration  of the device
	Pin struct {
		Required bool        `json:"required,omitempty"`
		Pattern  interface{} `json:"pattern,omitempty"`
	} `json:"pin,omitempty"`
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
}

func NewThing(id string, description []byte) (thing *Thing) {
	var th Thing
	th.ID = id
	err := json.Unmarshal(description, &th)
	if len(th.AtContext) == 0 {
		th.AtContext = []string{"https://webthings.io/schemas/"}
	}
	if err != nil {
		return nil
	}
	return &th
}

func (t *Thing) findProperty(propName string) (*Property, error) {
	prop, ok := t.Properties[propName]
	if !ok {
		return nil, fmt.Errorf("thing(%s) can not found properties(%s)", t.ID, propName)
	}
	return prop, nil
}

func (t *Thing) GetProperty(propName string) *Property {
	prop, ok := t.Properties[propName]
	if !ok {
		log.Debug("thing(%s) can not found properties(%s)", t.ID, propName)
		return nil
	}
	return prop
}

func (t *Thing) GetId() string {
	return t.ID
}

func (t *Thing) SetTitle(title string) {
	t.Title = title
}

func (t *Thing) GetTitle() string {
	return t.Title
}

func (t *Thing) AddAction(action *Action) error {
	return nil
}

func (t *Thing) SetSelectedCapability(selectedCapability string) {
	t.SelectedCapability = selectedCapability
}

func (t *Thing) SetConnected(connected bool) {
	t.Connected = connected

}

func (t *Thing) IsConnected() bool {
	return t.Connected
}

func (t *Thing) RemoveAction(a *Action) bool {
	_, ok := t.Actions[a.Name]
	return ok
}

//thing save to database must do this:
func (t *Thing) GetDescription() string {
	s, err := json.MarshalToString(t)
	if err != nil {
		return ""
	}
	return s
}

func (t *Thing) Save() (err error) {
	return database.SetSetting(t.ID, t.GetDescription())
}

func (t *Thing) Update(thing *Thing) {
	_ = database.UpdateThing(t.ID, thing.GetDescription())
}
