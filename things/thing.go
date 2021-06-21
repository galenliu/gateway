package things

import (
	"fmt"
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/util"
	core2 "github.com/galenliu/gateway/pkg/wot/definitions/core"
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	hypermedia_controls2 "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/galenliu/gateway/wot/definitions/core"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"strings"
)

type Thing struct {
	*core2.Thing
	//The configuration  of the device
	Pin                 *addon.PIN `json:"pin,omitempty"`
	CredentialsRequired bool       `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(description string) (thing *Thing, err error) {
	bt := []byte(description)
	id := json.Get(bt, "id").ToString()
	if id == "" {
		return nil, fmt.Errorf("invaild id")
	}
	title := json.Get(bt, "title").ToString()
	if title == "" {
		title = id
	}

	sl := strings.Split(id, "/")
	tid := sl[len(sl)-1]
	thingId := fmt.Sprintf("%s/%s", util.ThingsPath, tid)

	var atContext []string
	for _, c := range gjson.Get(description, "@context").Array() {
		atContext = append(atContext, c.String())
	}

	var atType []string
	for _, c := range gjson.Get(description, "@type").Array() {
		atType = append(atType, c.String())
	}
	if len(atType) < 1 {
		return nil, fmt.Errorf("new thing err: @type")
	}

	t := &Thing{Thing: &core2.Thing{
		AtContext: atContext,
		AtType:    atType,
		Title:     title,
		ID:        thingId,
	}}

	t.IconData = gjson.Get(description, "iconData").String()
	t.Connected = gjson.Get(description, "connected").Bool()

	if gjson.Get(description, "description").Exists() {
		t.Description = gjson.Get(description, "description").String()
	}

	if gjson.Get(description, "credentialsRequired").Exists() {
		t.CredentialsRequired = gjson.Get(description, "credentialsRequired").Bool()
	}

	if gjson.Get(description, "selectedCapability").Exists() {
		sc := gjson.Get(description, "selectedCapability").String()
		for _, s := range t.AtType {
			if s == sc {
				t.SelectedCapability = sc
			}
		}
	} else {
		if len(t.AtType) > 0 {
			t.SelectedCapability = t.AtType[0]
		}
	}

	if gjson.Get(description, "pin").Exists() {
		var pin addon.PIN
		pin.Required = gjson.Get(description, "pin.required").Bool()
		pin.Pattern = gjson.Get(description, "pin.pattern").Value()
		t.Pin = &pin
	}

	if gjson.Get(description, "properties").Exists() {
		var props = gjson.Get(description, "properties").Map()
		if len(props) > 0 {
			t.Properties = make(map[string]core2.PropertyAffordance)
			for name, data := range props {
				prop := core.NewPropertyFromString(data.String())
				if prop != nil {
					if prop.Forms == nil {
						prop.Forms = append(prop.Forms, hypermedia_controls2.Form{
							Href:        fmt.Sprintf("%s%s/%s", thingId, util.PropertiesPath, name),
							ContentType: data_schema2.ApplicationJson,
							Op:          []string{hypermedia_controls2.ReadProperty, hypermedia_controls2.WriteProperty},
						})
					}
					t.Properties[name] = prop
				}
			}
		}
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Op: []string{hypermedia_controls2.ReadallProperties, hypermedia_controls2.WriteAllProperties}, Href: thingId + util.PropertiesPath, ContentType: data_schema2.ApplicationJson})
	}

	if gjson.Get(description, "_actions").Exists() {
		var actions = gjson.Get(description, "_actions").Map()
		if len(actions) > 0 {
			t.Actions = make(map[string]*model.Action)
			for name, a := range actions {
				var action model.Action
				err := json.UnmarshalFromString(a.String(), &action)
				if err != nil {
					continue
				}
				action.ID = a.Get("id").String()
				if action.Forms == nil {
					action.Forms = append(action.Forms, hypermedia_controls2.Form{Href: fmt.Sprintf("%s/_actions/%s", t.ID, name)})
				}
				if action.Name == "" {
					action.Name = name
				}
				t.Actions[name] = &action
			}

			t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "_actions", Href: thingId + util.ActionsPath})
		}
	}

	if gjson.Get(description, "events").Exists() {
		var events = gjson.Get(description, "events").Map()
		if len(events) > 0 {
			t.Events = make(map[string]*model.Event)
			for name, e := range events {
				var event model.Event
				err := json.UnmarshalFromString(e.String(), &event)
				if err != nil {
					continue
				}
				if !e.Get("id").Exists() {
					continue
				}
				event.ID = e.Get("id").String()
				if event.Forms == nil {
					event.Forms = append(event.Forms, hypermedia_controls2.Form{Href: fmt.Sprintf("%s/events/%s", thingId, name)})
				}
				t.Events[name] = &event
			}
			t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "_actions", Href: t.ID + util.EventsPath})
		}
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "alternate", ContentType: "text/html", Href: t.ID})
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "alternate", Href: fmt.Sprintf("wss://localhost/%s", thingId)})
	}
	return t, nil
}
