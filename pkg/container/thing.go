package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

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
	IconData           string `json:"iconData,omitempty"`
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(description string) (thing *Thing, err error) {
	data := []byte(description)

	t := Thing{}
	t.Thing, err = wot.NewThingFromString(description)
	t.Thing.ID = hypermedia_controls.URI(fmt.Sprintf("%s/%s", constant.ThingsPath, thing.ID.GetID()))
	if err != nil {
		return nil, err
	}

	t.IconData = json.Get(data, "iconData").ToString()
	t.Connected = json.Get(data, "connected").ToBool()
	t.CredentialsRequired = json.Get(data, "credentialsRequired").ToBool()

	sc := json.Get(data, "selectedCapability").ToString()
	for _, s := range t.AtType {
		if s == sc {
			t.SelectedCapability = sc
			break
		}
		continue
	}
	if t.SelectedCapability == "" && len(t.AtType) > 0 {
		t.SelectedCapability = t.AtType[0]
	}

	var p PIN
	json.Get(data, "pin").ToVal(&p)
	if &p != nil {
		t.Pin = &p
	}
	return &t, nil
}
