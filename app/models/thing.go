package models

import (
	"fmt"
	"gateway/addons"
	"gateway/pkg/database"
)

type ThingInfo struct {
	ID          string   `json:"id"`
	AtContext   string   `json:"@context,omitempty" gorm:"@context"`
	AtType      []string `json:"@type" gorm:"@type"`
	Title       string   `json:"title,required"`
	Description string   `json:"description,omitempty"`

	Links    []string `json:"links"`
	BaseHref string   `json:"baseHref"`
	Pin      struct {
		Required bool   `json:"required"`
		Pattern  string `json:"pattern"`
	} `json:"pin,omitempty"`

	Href                string `json:"href"`
	CredentialsRequired bool   `json:"credentialsRequired"`
	SelectedCapability  string `json:"selectedCapability"`

	Properties map[string]*Property `json:"properties" gorm:"-"`
	Actions    map[string]*Action   `json:"actions" gorm:"-"`
	Events     map[string]*Event    `json:"events" gorm:"-"`
}

type Thing struct {
	*ThingInfo
	Properties []*Property `json:"-"`
	Actions    []*Action   `json:"-"`
	Events     []*Event    `json:"-"`

	connected bool
}

func NewThing(t *ThingInfo) *Thing {
	th := &Thing{ThingInfo: t}
	return th
}

func (t *Thing) SetTitle(title string) {
	t.Title = title
}

func (t *Thing) AddAction(action *Action) error {
	return nil
}

func (t *Thing) SetSelectedCapability(selectedCapability string) {
	t.SelectedCapability = selectedCapability
}

func deviceToThing(devices map[string]*addons.DeviceProxy) map[string]*ThingInfo {
	var thingsMap = make(map[string]*ThingInfo)
	for key, dev := range devices {
		thing := &ThingInfo{
			ID:                  dev.ID,
			AtContext:           dev.AtContext,
			AtType:              dev.AtType,
			Title:               dev.Title,
			Description:         dev.Description,
			Links:               nil,
			BaseHref:            "",
			Href:                "",
			CredentialsRequired: false,
			Properties:          nil,
		}

		var props = make(map[string]*Property)
		for _, p := range dev.Properties {
			thingProp := devPropToThingProp(p.Property, thing.ID)
			props[thingProp.Name] = thingProp
		}
		thing.Properties = props
		thingsMap[key] = thing
	}
	return thingsMap
}

func (t Thing) saveToDataBase() error {
	if len(t.ThingInfo.Properties) > 0 {
		for _, p := range t.ThingInfo.Properties {
			t.Properties = append(t.Properties, p)
		}
	}
	if len(t.ThingInfo.Actions) > 0 {
		for _, a := range t.ThingInfo.Actions {
			t.Actions = append(t.Actions, a)
		}
	}
	if len(t.ThingInfo.Events) > 0 {
		for _, e := range t.ThingInfo.Events {
			t.Events = append(t.Events, e)
		}
	}
	db, err := database.GetDB()

	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Thing{}, &Action{}, &Event{})
	if err != nil {
		return err
	}
	tx := db.Create(t)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetThingByIdFormDataBase(id string) (*Thing, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	var thing Thing
	tx := db.First(&thing, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected < 1 {
		return nil, fmt.Errorf("can not find thing(%s)", id)
	}
	thing.update()
	return &thing, nil
}

func GetThingsFormDataBase() (*[]Thing, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	var things []Thing
	tx := db.Find(&things)
	if tx.Error != nil {
		return nil, err
	}
	if tx.RowsAffected < 1 {
		return nil, fmt.Errorf("conut 0")
	}
	for _, t := range things {
		t.update()
	}

	return &things, nil
}

func (t *Thing) update() {
	t.ThingInfo.Properties = map[string]*Property{}
	t.ThingInfo.Actions = map[string]*Action{}
	t.ThingInfo.Events = map[string]*Event{}

	if len(t.Properties) > 0 {
		for _, p := range t.Properties {
			t.ThingInfo.Properties[p.Name] = p
		}
	}
	if len(t.Actions) > 0 {
		for _, a := range t.Actions {
			t.ThingInfo.Actions[a.Name] = a
		}
	}
	if len(t.Events) > 0 {
		for _, e := range t.Events {
			t.ThingInfo.Events[e.Name] = e
		}
	}
}
