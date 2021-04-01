package thing

import (
	"addon"
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
)

type Thing struct {
	AtContext   []string `json:"@context"`
	Title       string   `json:"title"`
	ID          string   `json:"id"`
	AtType      []string `json:"@type"`
	Description string   `json:"description,omitempty"`

	Properties map[string]*Property `json:"properties"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`
	Forms      []util.Form          `json:"forms,omitempty"`

	//The configuration  of the device
	Pin                 addon.PIN
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
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

func (t *Thing) setSelectedCapability(sel string) {
	t.SelectedCapability = sel
	_ = t.update
	bus.Publish(util.MODIFIED, t)
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

func (t *Thing) SetTitle(title string) string {
	t.Title = title
	_ = t.update
	bus.Publish(util.MODIFIED, t)
	return t.GetDescription()
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

func (t *Thing) save() (err error) {
	return database.SetSetting(t.ID, t.GetDescription())
}

func (t *Thing) update(thing *Thing) {
	_ = database.UpdateThing(t.ID, thing.GetDescription())
}
